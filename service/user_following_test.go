package service

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestIsFollowing(t *testing.T) {
	t.Run("should return false if no rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WillReturnError(sql.ErrNoRows)
		if isFollowing := userService.IsFollowing(testCtx, "follower-id", "following-id"); isFollowing {
			t.Fatal("isFollowing should return false")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should return true if return rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"following_id"}).
				AddRow("some-id")).
			RowsWillBeClosed()
		if isFollowing := userService.IsFollowing(testCtx, "follower-id", "following-id"); !isFollowing {
			t.Fatal("isFollowing should return true")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestFollowUser(t *testing.T) {
	t.Run("should return error if following the user itself", func(t *testing.T) {
		id := "my-id"
		ro := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "bio", "image"}).
			AddRow(id, "example@mail.com", "example", "password", "", "")
		mock.ExpectQuery("SELECT").WillReturnRows(ro).RowsWillBeClosed()

		_, err := userService.FollowUser(testCtx, id, "example")
		if err == nil {
			t.Fatal("expected return error")
		}

		if err.Code != http.StatusBadRequest {
			t.Fatalf("expected return code: %d but got: %d\n", http.StatusBadRequest, err.Code)
		}

		if err.Error != ErrSelfFollow {
			t.Fatalf("expected return error: %v but got: %v\n", ErrSelfFollow, err.Error)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestUnfollowUser(t *testing.T) {
	t.Run("should return error if unfollowing the user itself", func(t *testing.T) {
		id := "my-id"
		ro := sqlmock.
			NewRows([]string{"id", "email", "username", "password", "bio", "image"}).
			AddRow(id, "example@mail.com", "example", "password", "", "")
		mock.ExpectQuery("SELECT").WillReturnRows(ro).RowsWillBeClosed()

		_, err := userService.UnfollowUser(testCtx, id, "example")
		if err == nil {
			t.Fatal("should return error")
		}

		if err.Code != http.StatusBadRequest {
			t.Fatalf("should return code: %d but got: %d\n", http.StatusBadRequest, err.Code)
		}

		if err.Error != ErrSelfUnfollow {
			t.Fatalf("should return error: %v but got: %v\n", ErrSelfUnfollow, err.Error)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}
