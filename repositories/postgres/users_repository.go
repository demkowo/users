// repository/repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	model "github.com/demkowo/users/models"
	"github.com/demkowo/users/repositories/postgres/sqlc"

	"github.com/google/uuid"
)

type Users interface {
	Add(ctx context.Context, user model.User) (model.User, error)
	Find(ctx context.Context) ([]model.User, error)
	List(ctx context.Context, limit, offset int32) ([]model.User, error)
	GetImgByNickname(ctx context.Context, nickname string) (string, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	Update(ctx context.Context, user model.User) (model.User, error)
	UpdateImg(ctx context.Context, userID uuid.UUID, img string) (model.User, error)
	Delete(ctx context.Context, userID uuid.UUID) (model.User, error)
}

type users struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewUsers(db *sql.DB) Users {
	return &users{
		db: db,
		q:  sqlc.New(db),
	}
}

func (r *users) Add(ctx context.Context, user model.User) (model.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.User{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	qtx := r.q.WithTx(tx)

	u, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Nickname: user.Nickname,
		Img:      nullString(user.Img),
		Country:  nullString(user.Country),
		City:     nullString(user.City),
	})
	if err != nil {
		_ = tx.Rollback()
		return model.User{}, err
	}

	// Insert clubs and user_clubs links
	for _, c := range user.Clubs {
		cl, err := qtx.CreateClub(ctx, c.Name)
		if err != nil {
			_ = tx.Rollback()
			return model.User{}, err
		}
		if err := qtx.AddUserClub(ctx, sqlc.AddUserClubParams{
			UserID: u.ID,
			ClubID: cl.ID,
		}); err != nil {
			_ = tx.Rollback()
			return model.User{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return model.User{}, err
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return model.User{}, err
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) Find(ctx context.Context) ([]model.User, error) {
	us, err := r.q.FindUsers(ctx)
	if err != nil {
		return nil, err
	}
	return r.attachClubs(ctx, us)
}

func (r *users) List(ctx context.Context, limit, offset int32) ([]model.User, error) {
	us, err := r.q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return r.attachClubs(ctx, us)
}

func (r *users) GetImgByNickname(ctx context.Context, nickname string) (string, error) {
	img, err := r.q.GetUserImgByNickname(ctx, nickname)
	if err != nil {
		return "", err
	}
	return nullStringToString(img), nil
}

func (r *users) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	if nullBoolToBool(u.Deleted) {
		return model.User{}, errors.New("user not found or deleted")
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return model.User{}, err
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) Update(ctx context.Context, user model.User) (model.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.User{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	qtx := r.q.WithTx(tx)

	u, err := qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:      user.ID,
		Country: nullString(user.Country),
		City:    nullString(user.City),
	})
	if err != nil {
		_ = tx.Rollback()
		return model.User{}, err
	}

	// Delete old links
	if err := qtx.DeleteUserClubsByUserID(ctx, u.ID); err != nil {
		_ = tx.Rollback()
		return model.User{}, err
	}

	// Insert new clubs
	for _, c := range user.Clubs {
		cl, err := qtx.CreateClub(ctx, c.Name)
		if err != nil {
			_ = tx.Rollback()
			return model.User{}, err
		}
		if err := qtx.AddUserClub(ctx, sqlc.AddUserClubParams{
			UserID: u.ID,
			ClubID: cl.ID,
		}); err != nil {
			_ = tx.Rollback()
			return model.User{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return model.User{}, err
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return model.User{}, err
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) UpdateImg(ctx context.Context, userID uuid.UUID, img string) (model.User, error) {
	u, err := r.q.UpdateUserImg(ctx, sqlc.UpdateUserImgParams{
		ID:  userID,
		Img: nullString(img),
	})
	if err != nil {
		return model.User{}, err
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return model.User{}, err
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) Delete(ctx context.Context, userID uuid.UUID) (model.User, error) {
	u, err := r.q.SoftDeleteUser(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		// user is already soft-deleted, return user with empty clubs
		return toDomainUser(u, []sqlc.Club{}), nil
	}

	return toDomainUser(u, clubs), nil
}

// attachClubs fetches clubs for each user
func (r *users) attachClubs(ctx context.Context, us []sqlc.User) ([]model.User, error) {
	domainUsers := toDomainUsers(us)
	for i, u := range us {
		clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		domainUsers[i].Clubs = clubsToDomain(clubs)
	}
	return domainUsers, nil
}

// Helpers
func toDomainUser(u sqlc.User, clubs []sqlc.Club) model.User {
	return model.User{
		ID:       u.ID,
		Nickname: u.Nickname,
		Img:      nullStringToString(u.Img),
		Country:  nullStringToString(u.Country),
		City:     nullStringToString(u.City),
		Created:  nullTimeToTime(u.CreatedAt),
		Updated:  nullTimeToTime(u.UpdatedAt),
		Deleted:  nullBoolToBool(u.Deleted),
		Clubs:    clubsToDomain(clubs),
	}
}

func toDomainUsers(us []sqlc.User) []model.User {
	users := make([]model.User, len(us))
	for i, u := range us {
		users[i] = model.User{
			ID:       u.ID,
			Nickname: u.Nickname,
			Img:      nullStringToString(u.Img),
			Country:  nullStringToString(u.Country),
			City:     nullStringToString(u.City),
			Created:  nullTimeToTime(u.CreatedAt),
			Updated:  nullTimeToTime(u.UpdatedAt),
			Deleted:  nullBoolToBool(u.Deleted),
		}
	}
	return users
}

func clubsToDomain(cs []sqlc.Club) []model.Club {
	clubs := make([]model.Club, len(cs))
	for i, c := range cs {
		clubs[i] = model.Club{
			ID:   c.ID,
			Name: c.Name,
		}
	}
	return clubs
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullTimeToTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

func nullBoolToBool(nb sql.NullBool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return false
}
