package service

import (
	"context"
	"net/http"

	model "github.com/demkowo/users/models"
	"github.com/demkowo/utils/resp"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type UsersRepo interface {
	Add(ctx context.Context, user model.User) (model.User, *resp.Err)
	Find(ctx context.Context) ([]model.User, *resp.Err)
	List(ctx context.Context, limit, offset int32) ([]model.User, *resp.Err)
	GetImgByNickname(ctx context.Context, nickname string) (string, *resp.Err)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, *resp.Err)
	Update(ctx context.Context, user model.User) (model.User, *resp.Err)
	UpdateImg(ctx context.Context, userID uuid.UUID, img string) (model.User, *resp.Err)
	Delete(ctx context.Context, userID uuid.UUID) (model.User, *resp.Err)
}

type Users interface {
	Find(ctx context.Context) ([]*model.User, *resp.Err)
	List(ctx context.Context, limit, offset int32) ([]*model.User, *resp.Err)
	GetAvatarByNickname(ctx context.Context, nickname string) (string, *resp.Err)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, *resp.Err)
	Add(ctx context.Context, user *model.User) *resp.Err
	Update(ctx context.Context, user *model.User) *resp.Err
	UpdateImg(ctx context.Context, id uuid.UUID, path string) *resp.Err
	Delete(ctx context.Context, id string) *resp.Err
}

type users struct {
	repo UsersRepo
}

func NewUsers(repo UsersRepo) Users {
	log.Trace()
	return &users{
		repo: repo,
	}
}

func (s *users) Find(ctx context.Context) ([]*model.User, *resp.Err) {
	log.Trace()

	us, err := s.repo.Find(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.User, len(us))
	for i := range us {
		u := us[i]
		result[i] = &u
	}
	return result, nil
}

func (s *users) List(ctx context.Context, limit, offset int32) ([]*model.User, *resp.Err) {
	log.Trace()

	us, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*model.User, len(us))
	for i := range us {
		u := us[i]
		result[i] = &u
	}
	return result, nil
}

func (s *users) GetAvatarByNickname(ctx context.Context, nickname string) (string, *resp.Err) {
	log.Trace()

	avatar, err := s.repo.GetImgByNickname(ctx, nickname)
	if err != nil {
		return "", err
	}
	return avatar, nil
}

func (s *users) GetByID(ctx context.Context, id uuid.UUID) (*model.User, *resp.Err) {
	log.Trace()

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *users) Add(ctx context.Context, user *model.User) *resp.Err {
	log.Trace()

	_, err := s.repo.Add(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}

func (s *users) Update(ctx context.Context, user *model.User) *resp.Err {
	log.Trace()

	_, err := s.repo.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}

func (s *users) UpdateImg(ctx context.Context, id uuid.UUID, path string) *resp.Err {
	log.Trace()

	_, err := s.repo.UpdateImg(ctx, id, path)
	if err != nil {
		return err
	}

	return nil
}

func (s *users) Delete(ctx context.Context, id string) *resp.Err {
	log.Trace()

	if id == "" {
		return resp.Error(http.StatusBadRequest, "failed to delete user", []interface{}{"user id is required"})
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return resp.Error(http.StatusBadRequest, "failed to delete user", []interface{}{"invalid uuid format", err.Error()})
	}

	_, e := s.repo.Delete(ctx, uid)
	if e != nil {
		return e
	}

	return nil
}
