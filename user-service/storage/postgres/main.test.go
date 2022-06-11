package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/project1/user-service/config"
	"github.com/project1/user-service/pkg/db"
	"github.com/project1/user-service/pkg/logger"
)

var repo *userRepo

func TestMain(m testing.M) {
	cfg := config.Load()

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	repo = NewUserRepo(connDB)

	os.Exit(m.Run())
}
