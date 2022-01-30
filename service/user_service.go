package service

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
)

type UserService struct {
	userRepo   *repository.UserRepository
	followRepo *repository.FollowingRepository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		userRepo:   repo.UserRepo,
		followRepo: repo.FollowRepo,
	}
}

func (s *UserService) GetOneById(id string) (*model.User, *ServiceError) {
	u := &model.User{}
	if err := s.userRepo.FindOneById(id, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, errors.New("no user found"))
		}
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return u, nil
}

func (s *UserService) GetOne(d *dto.LoginUserDto) (*model.User, *ServiceError) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}

	if err := s.userRepo.FindOne(u); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, errors.New("no user found"))
		}
		return nil, CreateServiceError(http.StatusInternalServerError, err)
	}

	return u, nil
}

func (s *UserService) CreateOne(d *dto.RegisterUserDto) (*model.User, *ServiceError) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}

	if err := u.HashPassword(d.Password); err != nil {
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	if err := s.userRepo.InsertOne(u); err != nil {
		switch err.Error() {
		case ErrDuplicateEmail:
			return nil, CreateServiceError(http.StatusBadRequest, errors.New("email already exist"))
		case ErrDuplicateUsername:
			return nil, CreateServiceError(http.StatusBadRequest, errors.New("username already exist"))
		default:
			return nil, CreateServiceError(http.StatusInternalServerError, err)
		}

	}

	return u, nil
}

func (s *UserService) Update(d *dto.UpdateUserDto, uid string) error {
	u := &model.User{
		ID:       uid,
		Email:    d.Email,
		Username: d.Username,
		Bio:      d.Bio,
		Image:    d.Image,
	}

	if err := u.HashPassword(d.Password); err != nil {
		return conduit.ErrInternal
	}

	return s.userRepo.UpdateOne(u)
}

func (s *UserService) FollowUser(followerID, username string) *ServiceError {
	following, err := s.GetOne(&dto.LoginUserDto{Username: username})
	if err != nil {
		return err
	}

	if followerID == following.ID {
		return CreateServiceError(http.StatusBadRequest, errors.New("you cannot follow your self"))
	}

	if err := s.followRepo.Follow(followerID, following.ID); err != nil {
		switch err.Error() {
		case ErrDuplicateFollowing:
			return CreateServiceError(http.StatusBadRequest, errors.New("you are already follow this user"))
		default:
			return CreateServiceError(http.StatusInternalServerError, err)
		}
	}

	return nil
}

func (s *UserService) UnfollowUser(followerID, username string) *ServiceError {
	following, err := s.GetOne(&dto.LoginUserDto{Username: username})
	if err != nil {
		return err
	}

	if followerID == following.ID {
		return CreateServiceError(http.StatusBadRequest, errors.New("you cannot unfollow your self"))
	}

	if err := s.followRepo.Unfollow(followerID, following.ID); err != nil {
		return CreateServiceError(http.StatusInternalServerError, err)
	}

	return nil
}

func (s *UserService) IsFollowing(followerID, followingID string) bool {
	return s.followRepo.IsFollowing(followerID, followingID)
}
