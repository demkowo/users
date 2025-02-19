package app

import (
	"database/sql"
	"os"

	"github.com/demkowo/users/internal/config"
	handler "github.com/demkowo/users/internal/handlers/gin"
	"github.com/demkowo/users/internal/repositories/postgres"
	service "github.com/demkowo/users/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const portNumber = ":5000"

var (
	conf         = config.Values.Get()
	dbConnection = os.Getenv("DB_CONNECTION")
	router       = gin.Default()
)

func init() {
	conf.UseCache = false
	conf.InProduction = false
	config.Values.Set(*conf)
}

func Start() {
	log.Trace()

	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	usersRepo := postgres.NewUsers(db)
	usersService := service.NewUsers(usersRepo)
	usersHandler := handler.NewUser(usersService)
	addUserRoutes(usersHandler)

	go pbServerStart(usersService)

	log.Infof("Starting server on %s", portNumber)
	if err := router.Run(portNumber); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
