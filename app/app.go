package app

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/demkowo/users/config"
	handler "github.com/demkowo/users/handlers"
	"github.com/demkowo/users/repositories/postgres"
	service "github.com/demkowo/users/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const portNumber = ":5001"

var (
	conf         = config.Values.Get()
	dbConnection = os.Getenv("DB_CONNECTION")
	router       = gin.Default()
)

func init() {
	// Basic config for logging and templates if needed
	fmt.Println("=== Basic logger configuration ===")
	conf.UseCache = false
	conf.InProduction = false
	config.Values.Set(*conf)
}

func Start() {
	log.Trace()

	// Open DB connection
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	usersRepo := postgres.NewUsers(db)
	usersService := service.NewUsers(usersRepo)
	usersHandler := handler.NewUser(usersService)
	addUserRoutes(usersHandler)

	log.Infof("Starting server on %s", portNumber)
	if err := router.Run(portNumber); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
