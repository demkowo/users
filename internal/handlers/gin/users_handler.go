package handler

import (
	"net/http"
	"strconv"

	model "github.com/demkowo/users/internal/models"
	service "github.com/demkowo/users/internal/services"
	"github.com/demkowo/utils/helper"
	"github.com/demkowo/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	help = helper.NewHelper()
)

type Users interface {
	Add(*gin.Context)
	Delete(*gin.Context)
	Find(*gin.Context)
	GetAvatarByNickname(*gin.Context)
	GetById(*gin.Context)
	List(*gin.Context)
	Update(*gin.Context)
	UpdateImg(*gin.Context)
}

type users struct {
	service service.Users
}

func NewUser(service service.Users) Users {
	log.Trace()
	return &users{
		service: service,
	}
}

func (h *users) Add(c *gin.Context) {
	log.Trace()

	ctx := c.Request.Context()
	var input struct {
		Nickname string   `json:"nickname"`
		Img      string   `json:"img"`
		Country  string   `json:"country"`
		City     string   `json:"city"`
		Clubs    []string `json:"clubs"`
	}

	if !help.BindJSON(c, &input) {
		return
	}

	var clubs []model.Club
	for _, club := range input.Clubs {
		clubs = append(clubs, model.Club{
			ID:   uuid.New(),
			Name: club,
		})
	}

	user := &model.User{
		ID:       uuid.New(),
		Nickname: input.Nickname,
		Img:      input.Img,
		Country:  input.Country,
		City:     input.City,
		Clubs:    clubs,
	}

	if err := h.service.Add(ctx, user); err != nil {
		log.Errorf("Failed to create user: %v", err)
		c.JSON(err.JSON())
		return
	}

	c.JSON(resp.New(http.StatusOK, "user added succesfully", []interface{}{user}).JSON())
}

func (h *users) Delete(c *gin.Context) {
	log.Trace()

	ctx := c.Request.Context()
	userID := c.Param("user_id")
	if err := h.service.Delete(ctx, userID); err != nil {
		log.Errorf("Failed to delete user: %v", err)
		c.JSON(err.JSON())
		return
	}

	c.JSON(resp.New(http.StatusOK, "user deleted successfully", nil).JSON())
}

func (h *users) Find(c *gin.Context) {
	log.Trace()

	ctx := c.Request.Context()
	users, err := h.service.Find(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(err.JSON())
		return
	}

	c.JSON(resp.New(http.StatusOK, "users found successfully", []interface{}{users}).JSON())
}

func (h *users) GetAvatarByNickname(c *gin.Context) {
	log.Trace()

	ctx := c.Request.Context()
	nickname := c.Param("nickname")
	avatar, err := h.service.GetAvatarByNickname(ctx, nickname)
	if err != nil {
		log.Error(err)
		c.JSON(err.JSON())
		return
	}

	c.JSON(resp.New(http.StatusOK, "avatar fetched successfully", []interface{}{avatar}).JSON())
}

func (h *users) GetById(c *gin.Context) {
	log.Trace()

	ctx := c.Request.Context()
	idStr := c.Param("user_id")
	if idStr == "" {
		log.Error("empty user ID")
		c.JSON(resp.Error(http.StatusBadRequest, "failed to get user", []interface{}{"user id can't be empty"}).JSON())
		return
	}

	var id uuid.UUID
	if !help.ParseUUID(c, "user_id", idStr, &id) {
		return
	}

	user, e := h.service.GetByID(ctx, id)
	if e != nil {
		log.Error(e)
		c.JSON(e.JSON())
		return
	}

	c.JSON(resp.New(http.StatusOK, "user fetched successfully", []interface{}{user}).JSON())
}

func (h *users) List(c *gin.Context) {
	log.Trace()

	ctx := c.Request.Context()
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var (
		limit  int32 = 10
		offset int32 = 0
	)

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = int32(l)
		}
	}

	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil && o >= 0 {
			offset = int32(o)
		}
	}

	users, e := h.service.List(ctx, limit, offset)
	if e != nil {
		log.Error(e)
		c.JSON(e.JSON())
		return
	}

	c.JSON(resp.New(http.StatusOK, "users fetched successfully", []interface{}{users}).JSON())
}

func (h *users) Update(c *gin.Context) {
	log.Trace()

	ctx := c.Request.Context()

	var id uuid.UUID
	if !help.ParseUUID(c, "user_id", c.Param("user_id"), &id) {
		return
	}

	var input struct {
		Country string   `json:"country"`
		City    string   `json:"city"`
		Clubs   []string `json:"clubs"`
	}

	if !help.BindJSON(c, &input) {
		return
	}

	var clubs []model.Club
	for _, club := range input.Clubs {
		clubs = append(clubs, model.Club{
			ID:   uuid.New(),
			Name: club,
		})
	}

	user := &model.User{
		ID:      id,
		Country: input.Country,
		City:    input.City,
		Clubs:   clubs,
	}

	if err := h.service.Update(ctx, user); err != nil {
		log.Errorf("Failed to update user: %v", err)
		c.JSON(err.JSON())
		return
	}

	c.JSON(resp.New(http.StatusOK, "user updated succesfully", nil).JSON())
}

func (h *users) UpdateImg(c *gin.Context) {
	log.Trace()

	ctx := c.Request.Context()

	var id uuid.UUID
	if !help.ParseUUID(c, "user_id", c.Param("user_id"), &id) {
		return
	}

	var input struct {
		Img string `json:"img" binding:"required"`
	}

	if !help.BindJSON(c, &input) {
		return
	}

	if err := h.service.UpdateImg(ctx, id, input.Img); err != nil {
		log.Errorf("Failed to update user image: %v", err)
		c.JSON(err.JSON())
		return
	}

	c.JSON(resp.New(http.StatusOK, "user image updated succesfully", nil).JSON())
}
