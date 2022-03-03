package service

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestIsFollowing(t *testing.T) {
	t.Run("should return false if no rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WillReturnError(sql.ErrNoRows)
		assert.False(t, userService.IsFollowing(testCtx, "follower-id", "following-id"), "should return false")
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("should return true if return rows", func(t *testing.T) {
		ro := sqlmock.
			NewRows([]string{"following_id"}).
			AddRow("some-id")
		mock.ExpectQuery("SELECT").
			WillReturnRows(ro).
			RowsWillBeClosed()
		assert.True(t, userService.IsFollowing(testCtx, "follower-id", "following-id"), "should return true")
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestFollowUser(t *testing.T) {
	as := assert.New(t)
	id := "my-id"
	ro := sqlmock.
		NewRows([]string{"id", "email", "username", "password", "bio", "image"}).
		AddRow(id, "example@mail.com", "example", "password", "", "")
	mock.ExpectQuery("SELECT").WillReturnRows(ro).RowsWillBeClosed()

	_, err := userService.FollowUser(testCtx, id, "example")
	as.NotNil(err, "error should not be nil")
	as.Equal(err.Code, http.StatusBadRequest)
	as.Equal(err.Error, ErrSelfFollow)
	as.Nil(mock.ExpectationsWereMet())
}

func TestUnfollowUser(t *testing.T) {
	as := assert.New(t)
	id := "my-id"
	ro := sqlmock.
		NewRows([]string{"id", "email", "username", "password", "bio", "image"}).
		AddRow(id, "example@mail.com", "example", "password", "", "")
	mock.ExpectQuery("SELECT").WillReturnRows(ro).RowsWillBeClosed()

	_, err := userService.UnfollowUser(testCtx, id, "example")
	as.NotNil(err, "error should not be nil")
	as.Equal(err.Code, http.StatusBadRequest)
	as.Equal(err.Error, ErrSelfUnfollow)
	as.Nil(mock.ExpectationsWereMet())
}
