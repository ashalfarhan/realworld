package service

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashalfarhan/realworld/api/dto"
)

func TestCreateArticle(t *testing.T) {
	d := &dto.CreateArticleDto{
		Article: &dto.CreateArticleFields{
			Title:       "Example title",
			Description: "Example description",
			Body:        "Example body",
			TagList: []string{
				"reactjs",
				"typescript",
			},
		},
	}

	row := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("asd", time.Now(), time.Now())

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO articles").
		ExpectQuery().
		WithArgs(sqlmock.AnyArg(), d.Article.Title, d.Article.Description, d.Article.Body, "author-id").
		WillReturnRows(row)
	mock.ExpectCommit()

	for _, tag := range d.Article.TagList {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO article_tags").
			WithArgs(sqlmock.AnyArg(), tag).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	authorRow := sqlmock.
		NewRows([]string{"id", "email", "username", "bio", "image", "created_at", "updated_at"}).
		AddRow("author-id", "example@mail.com", "example", "", "", time.Now(), time.Now())

	mock.ExpectQuery("SELECT").WillReturnRows(authorRow).RowsWillBeClosed()

	_, err := articleService.CreateArticle(testCtx, d, "author-id")
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

}
