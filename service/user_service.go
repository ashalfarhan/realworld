package service

import (
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

func (s *UserService) GetOneById(id string) (*model.User, *ServiceError) {
	u := &model.User{}
	if err := s.userRepo.FindOneById(id, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, errors.New("no user found"))
		}

		s.logger.Printf("Cannot FindOneById for %s, Reason: %v", id, err)
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

		s.logger.Printf("Cannot FindOne for %#v, Reason: %v", d, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return u, nil
}

func (s *UserService) CreateOne(d *dto.RegisterUserDto) (*model.User, *ServiceError) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}

	u.Password = s.HashPassword(d.Password)

	if err := s.userRepo.InsertOne(u); err != nil {
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

func (s *UserService) Update(d *dto.UpdateUserDto, uid string) *ServiceError {
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

	if err := s.userRepo.UpdateOne(u); err != nil {
		switch err.Error() {
		case ErrDuplicateEmail:
			return CreateServiceError(http.StatusBadRequest, errors.New("email already exist"))
		case ErrDuplicateUsername:
			return CreateServiceError(http.StatusBadRequest, errors.New("username already exist"))
		default:
			s.logger.Printf("Cannot UpdateOne payload:%#v args:%#v, Reason: %v", *d, *u, err)
			return CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return nil
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
			s.logger.Printf("Cannot FollowUser followerID:%s following.ID:%s, Reason: %v", followerID, following.ID, err)
			return CreateServiceError(http.StatusInternalServerError, nil)
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
		s.logger.Printf("Cannot UnfollowUser followerID:%s following.ID:%s, Reason: %v", followerID, following.ID, err)
		return CreateServiceError(http.StatusInternalServerError, err)
	}

	return nil
}

func (s *UserService) IsFollowing(followerID, followingID string) bool {
	return s.followRepo.IsFollowing(followerID, followingID)
}

func (s *UserService) HashPassword(p string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Printf("Cannot HashPassword, Reason: %v", err)
	}

	return string(hashed)
}
