package postgres

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	model "github.com/demkowo/users/models"
	"github.com/demkowo/users/repositories/postgres/sqlc"
	"github.com/demkowo/utils/resp"

	"github.com/google/uuid"
)

type Users interface {
	Add(ctx context.Context, user model.User) (model.User, *resp.Err)
	Find(ctx context.Context) ([]model.User, *resp.Err)
	List(ctx context.Context, limit, offset int32) ([]model.User, *resp.Err)
	GetImgByNickname(ctx context.Context, nickname string) (string, *resp.Err)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, *resp.Err)
	Update(ctx context.Context, user model.User) (model.User, *resp.Err)
	UpdateImg(ctx context.Context, userID uuid.UUID, img string) (model.User, *resp.Err)
	Delete(ctx context.Context, userID uuid.UUID) (model.User, *resp.Err)
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

func (r *users) Add(ctx context.Context, user model.User) (model.User, *resp.Err) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to add user", []interface{}{err.Error()})
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
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to add user", []interface{}{err.Error()})
	}

	for _, c := range user.Clubs {
		cl, err := qtx.CreateClub(ctx, c.Name)
		if err != nil {
			_ = tx.Rollback()
			return model.User{}, resp.Error(http.StatusInternalServerError, "failed to add user", []interface{}{err.Error()})
		}
		if err := qtx.AddUserClub(ctx, sqlc.AddUserClubParams{
			UserID: u.ID,
			ClubID: cl.ID,
		}); err != nil {
			_ = tx.Rollback()
			return model.User{}, resp.Error(http.StatusInternalServerError, "failed to add user", []interface{}{err.Error()})
		}
	}

	if err := tx.Commit(); err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to add user", []interface{}{err.Error()})
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to add user", []interface{}{err.Error()})
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) Find(ctx context.Context) ([]model.User, *resp.Err) {
	us, err := r.q.FindUsers(ctx)
	if err != nil {
		return nil, resp.Error(http.StatusInternalServerError, "failed to find users", []interface{}{err.Error()})
	}
	return r.attachClubs(ctx, us)
}

func (r *users) List(ctx context.Context, limit, offset int32) ([]model.User, *resp.Err) {
	us, err := r.q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, resp.Error(http.StatusInternalServerError, "failed to list users", []interface{}{err.Error()})
	}
	return r.attachClubs(ctx, us)
}

func (r *users) GetImgByNickname(ctx context.Context, nickname string) (string, *resp.Err) {
	img, err := r.q.GetUserImgByNickname(ctx, nickname)
	if err != nil {
		return "", resp.Error(http.StatusInternalServerError, "failed to get users image", []interface{}{err.Error()})
	}
	return nullStringToString(img), nil
}

func (r *users) GetByID(ctx context.Context, id uuid.UUID) (model.User, *resp.Err) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to get user", []interface{}{err.Error()})
	}
	if nullBoolToBool(u.Deleted) {
		return model.User{}, resp.Error(http.StatusInternalServerError, "user not found", nil)
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to get user", []interface{}{err.Error()})
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) Update(ctx context.Context, user model.User) (model.User, *resp.Err) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update user", []interface{}{err.Error()})
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
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update user", []interface{}{err.Error()})
	}

	if err := qtx.DeleteUserClubsByUserID(ctx, u.ID); err != nil {
		_ = tx.Rollback()
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update user", []interface{}{err.Error()})
	}

	for _, c := range user.Clubs {
		cl, err := qtx.CreateClub(ctx, c.Name)
		if err != nil {
			_ = tx.Rollback()
			return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update user", []interface{}{err.Error()})
		}
		if err := qtx.AddUserClub(ctx, sqlc.AddUserClubParams{
			UserID: u.ID,
			ClubID: cl.ID,
		}); err != nil {
			_ = tx.Rollback()
			return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update user", []interface{}{err.Error()})
		}
	}

	if err := tx.Commit(); err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update user", []interface{}{err.Error()})
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update user", []interface{}{err.Error()})
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) UpdateImg(ctx context.Context, userID uuid.UUID, img string) (model.User, *resp.Err) {
	u, err := r.q.UpdateUserImg(ctx, sqlc.UpdateUserImgParams{
		ID:  userID,
		Img: nullString(img),
	})
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update users image", []interface{}{err.Error()})
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to update users image", []interface{}{err.Error()})
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) Delete(ctx context.Context, userID uuid.UUID) (model.User, *resp.Err) {
	u, err := r.q.SoftDeleteUser(ctx, userID)
	if err != nil {
		return model.User{}, resp.Error(http.StatusInternalServerError, "failed to delete user", []interface{}{err.Error()})
	}

	clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
	if err != nil {
		return toDomainUser(u, []sqlc.Club{}), nil
	}

	return toDomainUser(u, clubs), nil
}

func (r *users) attachClubs(ctx context.Context, us []sqlc.User) ([]model.User, *resp.Err) {
	domainUsers := toDomainUsers(us)
	for i, u := range us {
		clubs, err := r.q.GetClubsByUserID(ctx, u.ID)
		if err != nil {
			return nil, resp.Error(http.StatusInternalServerError, "failed to attach clubs", []interface{}{err.Error()})
		}
		domainUsers[i].Clubs = clubsToDomain(clubs)
	}
	return domainUsers, nil
}

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
