package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/project1/post-service/config"
	"github.com/project1/post-service/pkg/db"
	"github.com/project1/post-service/pkg/logger"
)

var repo *postRepo

func TestMain(m testing.M) {
	cfg := config.Load()

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	repo = NewPostRepo(connDB)

	os.Exit(m.Run())
}
