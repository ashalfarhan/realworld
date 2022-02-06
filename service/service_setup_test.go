package service

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/jmoiron/sqlx"
)

var (
	mock           sqlmock.Sqlmock
	userService    *UserService
	articleService *ArticleService
	testCtx        = context.TODO()
)

func TestMain(m *testing.M) {
	var mockDb *sql.DB
	var err error

	mockDb, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer mockDb.Close()

	db := sqlx.NewDb(mockDb, "sqlmock")
	defer db.Close()

	repo := repository.InitRepository(db)
	userService = NewUserService(repo)
	articleService = NewArticleService(repo)

	os.Exit(m.Run())
}
