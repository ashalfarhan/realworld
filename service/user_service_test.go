package service

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/jmoiron/sqlx"
)

var (
	mock sqlmock.Sqlmock
	err  error
	us   *UserService
	ctx  context.Context
)

func TestMain(m *testing.M) {
	var mockDb *sql.DB
	mockDb, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer mockDb.Close()

	db := sqlx.NewDb(mockDb, "sqlmock")
	defer db.Close()

	repo := repository.InitRepository(db)
	us = NewUserService(repo)
	ctx = context.TODO()

	os.Exit(m.Run())
}

func TestRegisterUser(t *testing.T) {
	d := &RegisterArgs{
		Email:    "asd",
		Username: "asd",
		Password: "asd",
	}

	mock.ExpectBegin()

	ret := mock.NewRows([]string{"id", "bio", "image"}).AddRow("asdasd", "", "")
	mock.ExpectPrepare("INSERT INTO users").
		ExpectQuery().
		WithArgs(d.Email, d.Username, sqlmock.AnyArg()).
		WillReturnRows(ret).
		RowsWillBeClosed()

	mock.ExpectCommit()

	u, err := us.Register(context.TODO(), d)
	if err != nil {
		t.Fatal(err)
	}

	if u.Password == d.Password {
		t.Fatalf("expected password to be hashed, but got %s", u.Password)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateUser(t *testing.T) {
	oldUname, oldMail := "user1", "user1@mail.com"
	newUname, newMail, newPass := "user01", "", "new-password"
	d := &dto.UpdateUserDto{
		User: &dto.UpdateUserFields{
			Email:    newMail,
			Username: newUname,
			Password: newPass,
		},
	}

	findRes := sqlmock.
		NewRows([]string{"id", "email", "username", "bio", "image", "created_at", "updated_at"}).
		AddRow("uuid", oldMail, oldUname, "", "", time.Now(), time.Now())

	mock.ExpectQuery("SELECT").WillReturnRows(findRes).RowsWillBeClosed()

	mock.ExpectBegin()
	mock.ExpectPrepare("UPDATE users SET").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	u, err := us.Update(context.TODO(), d, "")
	if err != nil {
		t.Fatal(err)
	}

	if u.Username != newUname {
		t.Fatalf("new username should be set, current: %s, newUname: %s", u.Username, newUname)
	}

	if u.Email == newMail {
		t.Fatal("email should not updated because empty, but changed")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
