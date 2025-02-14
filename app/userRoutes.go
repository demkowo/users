package app

import (
	log "github.com/sirupsen/logrus"

	handler "github.com/demkowo/users/handlers"
)

func addUserRoutes(h handler.Users) {
	log.Println("--- Setting User Routes ---")

	router.POST("/api/v1/users/add", h.Add)
	router.PUT("/api/v1/users/edit/:user_id", h.Update)
	router.PUT("/api/v1/users/edit-img/:user_id", h.UpdateImg)
	router.DELETE("/api/v1/users/delete/:user_id", h.Delete)
	router.GET("/api/v1/users/get/:user_id", h.GetById)
	router.GET("/api/v1/users/get-avatar/:nickname", h.GetAvatarByNickname)
	router.GET("/api/v1/users/find", h.Find)
	router.GET("/api/v1/users/list", h.List)
}
