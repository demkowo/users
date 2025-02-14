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
	Add(ctx context.Context, user model.User) (model.User, error)
	Find(ctx context.Context) ([]model.User, error)
	List(ctx context.Context, limit, offset int32) ([]model.User, error)
	GetImgByNickname(ctx context.Context, nickname string) (string, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	Update(ctx context.Context, user model.User) (model.User, error)
	UpdateImg(ctx context.Context, userID uuid.UUID, img string) (model.User, error)
	Delete(ctx context.Context, userID uuid.UUID) (model.User, error)
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
		return nil, resp.Error(http.StatusInternalServerError, "failed to find users", []interface{}{err.Error()})
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
		return nil, resp.Error(http.StatusInternalServerError, "failed to list users", []interface{}{err.Error()})
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
		return "", resp.Error(http.StatusInternalServerError, "failed to get avatar", []interface{}{err.Error()})
	}
	return avatar, nil
}

func (s *users) GetByID(ctx context.Context, id uuid.UUID) (*model.User, *resp.Err) {
	log.Trace()

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, resp.Error(http.StatusInternalServerError, "failed to get user", []interface{}{err.Error()})
	}

	return &u, nil
}

func (s *users) Add(ctx context.Context, user *model.User) *resp.Err {
	log.Trace()

	_, err := s.repo.Add(ctx, *user)
	if err != nil {
		return resp.Error(http.StatusInternalServerError, "failed to add user", []interface{}{err.Error()})
	}

	return nil
}

func (s *users) Update(ctx context.Context, user *model.User) *resp.Err {
	log.Trace()

	_, err := s.repo.Update(ctx, *user)
	if err != nil {
		return resp.Error(http.StatusInternalServerError, "failed to update user", []interface{}{err.Error()})
	}

	return nil
}

func (s *users) UpdateImg(ctx context.Context, id uuid.UUID, path string) *resp.Err {
	log.Trace()

	_, err := s.repo.UpdateImg(ctx, id, path)
	if err != nil {
		return resp.Error(http.StatusInternalServerError, "failed to update image", []interface{}{err.Error()})
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
		return resp.Error(http.StatusBadRequest, "failed to delete user", []interface{}{e.Error()})
	}

	return nil
}
