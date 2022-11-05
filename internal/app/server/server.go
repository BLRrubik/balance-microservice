package server

import (
	"balance-microservice/internal/app/config"
	"balance-microservice/internal/app/store"
	sqlstore "balance-microservice/internal/app/store/postgresImpl"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type ErrorType struct {
	Message string `json:"message"`
}

type server struct {
	router *gin.Engine
	store  store.Store
}

func newServer(store *sqlstore.Store) *server {
	var s = &server{
		router: gin.Default(),
		store:  store,
	}
	s.configureRouter()

	return s
}

func Start(config *config.Config) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		config.DatabaseProps.Host, config.DatabaseProps.Port,
		config.DatabaseProps.User, config.DatabaseProps.Password,
		config.DatabaseProps.DBName, config.DatabaseProps.SSLMode)

	db, err := newDB(psqlInfo)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer db.Close()

	store := sqlstore.NewStore(db)

	srv := newServer(store)

	return srv.router.Run(":" + config.ServerProps.Port)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
