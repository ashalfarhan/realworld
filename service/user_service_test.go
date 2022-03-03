package service_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashalfarhan/realworld/api/dto"
	. "github.com/ashalfarhan/realworld/service"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	as := assert.New(t)
	d := &RegisterArgs{
		Email:    "asd",
		Username: "asd",
		Password: "asd",
	}

	prepareMockRegisterUser(d)
	u, err := userService.Register(testCtx, d)
	as.Nil(err, "Register should not return error")
	as.NotEqual(u.Password, d.Password, "Result password should be hased")
	as.Nil(mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	as := assert.New(t)
	oldUname, oldMail := "user1", "user1@mail.com"
	newUname, newMail, newPass := "user01", "", "new-password"
	d := &dto.UpdateUserDto{
		User: &dto.UpdateUserFields{
			Email:    newMail,
			Username: newUname,
			Password: newPass,
		},
	}

	prepareMockUpdateUser(oldMail, oldUname)
	u, err := userService.Update(testCtx, d, "")
	as.Nil(err, "Update user should not return error")
	as.Equal(u.Username, newUname, "New username should be changed")
	as.NotEqual(u.Email, newMail, "New mail should not be changed")
	as.Nil(mock.ExpectationsWereMet())
}

func prepareMockRegisterUser(d *RegisterArgs) {
	mock.ExpectBegin()
	ret := mock.
		NewRows([]string{"id", "bio", "image"}).
		AddRow("asdasd", "", "")
	mock.ExpectPrepare("INSERT INTO users").
		ExpectQuery().
		WithArgs(d.Email, d.Username, sqlmock.AnyArg()).
		WillReturnRows(ret).RowsWillBeClosed()
	mock.ExpectCommit()
}

func prepareMockUpdateUser(oldMail, oldUname string) {
	findRes := sqlmock.
		NewRows([]string{"id", "email", "username", "bio", "image", "created_at", "updated_at"}).
		AddRow("uuid", oldMail, oldUname, "", "", time.Now(), time.Now())
	mock.ExpectQuery("SELECT").
		WillReturnRows(findRes).RowsWillBeClosed()
	mock.ExpectBegin()
	mock.ExpectPrepare("UPDATE users SET").
		ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}
