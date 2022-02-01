package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   *repository.UserRepository
	followRepo *repository.FollowingRepository
	logger     *log.Logger
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		userRepo:   repo.UserRepo,
		followRepo: repo.FollowRepo,
		logger:     conduit.NewLogger("user-service"),
	}
}

func (s *UserService) GetOneById(ctx context.Context, id string) (*model.User, *ServiceError) {
	u := &model.User{}
	if err := s.userRepo.FindOneById(ctx, id, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, errors.New("no user found"))
		}

		s.logger.Printf("Cannot FindOneById for %s, Reason: %v", id, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return u, nil
}

func (s *UserService) GetOne(ctx context.Context, d *dto.LoginUserDto) (*model.User, *ServiceError) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}

	if err := s.userRepo.FindOne(ctx, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, errors.New("no user found"))
		}

		s.logger.Printf("Cannot FindOne for %#v, Reason: %v", d, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return u, nil
}

func (s *UserService) CreateOne(ctx context.Context, d *dto.RegisterUserDto) (*model.User, *ServiceError) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}

	u.Password = s.HashPassword(d.Password)

	if err := s.userRepo.InsertOne(ctx, u); err != nil {
		switch err.Error() {
		case ErrDuplicateEmail:
			return nil, CreateServiceError(http.StatusBadRequest, errors.New("email already exist"))
		case ErrDuplicateUsername:
			return nil, CreateServiceError(http.StatusBadRequest, errors.New("username already exist"))
		default:
			s.logger.Printf("Cannot InsertOne for %#v, Reason: %v", d, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}

	}

	return u, nil
}

func (s *UserService) Update(ctx context.Context, d *dto.UpdateUserDto, uid string) *ServiceError {
	u := &conduit.UpdateUserArgs{
		ID:    uid,
		Bio:   d.Bio,
		Image: d.Image,
	}

	if len(d.Email) != 0 {
		u.Email = &d.Email
	}
	if len(d.Username) != 0 {
		u.Username = &d.Username
	}
	if len(d.Password) != 0 {
		hashed := s.HashPassword(d.Password)
		u.Password = &hashed
	}

	if d.Bio.Set {
		u.Bio = d.Bio
	}

	if d.Image.Set {
		u.Image = d.Image
	}

	if err := s.userRepo.UpdateOne(ctx, u); err != nil {
		switch err.Error() {
		case ErrDuplicateEmail:
			return CreateServiceError(http.StatusBadRequest, errors.New("email already exist"))
		case ErrDuplicateUsername:
			return CreateServiceError(http.StatusBadRequest, errors.New("username already exist"))
		default:
			s.logger.Printf("Cannot UpdateOne payload:%#v args:%#v, Reason: %v", d, u, err)
			return CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return nil
}

func (s *UserService) FollowUser(ctx context.Context, followerID, username string) (*model.User, *ServiceError) {
	following, err := s.GetOne(ctx, &dto.LoginUserDto{Username: username})
	if err != nil {
		return nil, err
	}

	if followerID == following.ID {
		return nil, CreateServiceError(http.StatusBadRequest, errors.New("you cannot follow your self"))
	}

	if err := s.followRepo.Follow(ctx, followerID, following.ID); err != nil {
		switch err.Error() {
		case ErrDuplicateFollowing:
			return nil, CreateServiceError(http.StatusBadRequest, errors.New("you are already follow this user"))
		default:
			s.logger.Printf("Cannot FollowUser followerID:%s following.ID:%s, Reason: %v", followerID, following.ID, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return following, nil
}

func (s *UserService) UnfollowUser(ctx context.Context, followerID, username string) (*model.User, *ServiceError) {
	following, err := s.GetOne(ctx, &dto.LoginUserDto{Username: username})
	if err != nil {
		return nil, err
	}

	if followerID == following.ID {
		return nil, CreateServiceError(http.StatusBadRequest, errors.New("you cannot unfollow your self"))
	}

	if err := s.followRepo.Unfollow(ctx, followerID, following.ID); err != nil {
		s.logger.Printf("Cannot UnfollowUser followerID:%s following.ID:%s, Reason: %v", followerID, following.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, err)
	}

	return following, nil
}

func (s *UserService) IsFollowing(ctx context.Context, followerID, followingID string) bool {
	return s.followRepo.IsFollowing(ctx, followerID, followingID)
}

func (s *UserService) HashPassword(p string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Printf("Cannot HashPassword, Reason: %v", err)
	}

	return string(hashed)
}
