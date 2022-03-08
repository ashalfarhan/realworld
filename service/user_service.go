package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   repository.UserRepository
	followRepo repository.FollowingRepository
	logger     *logrus.Entry
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		repo.UserRepo,
		repo.FollowRepo,
		conduit.NewLogger("service", "UserService"),
	}
}

func (s *UserService) GetOneById(ctx context.Context, id string) (*model.User, *ServiceError) {
	u, err := s.userRepo.FindOneByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, ErrNoUserFound)
		}
		s.logger.Errorf("Cannot FindOneById for %s, Reason: %v", id, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return u, nil
}

func (s *UserService) GetOne(ctx context.Context, d *repository.FindOneUserFilter) (*model.User, *ServiceError) {
	u, err := s.userRepo.FindOne(ctx, d)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, ErrNoUserFound)
		}
		s.logger.Errorf("Cannot FindOne for %+v, Reason: %v", d, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return u, nil
}

func (s *UserService) Insert(ctx context.Context, d *dto.RegisterUserFields) (*model.User, *ServiceError) {
	d.Password = s.HashPassword(d.Password)
	u, err := s.userRepo.InsertOne(ctx, d)
	if err != nil {
		switch err.Error() {
		case repository.ErrDuplicateEmail:
			return nil, CreateServiceError(http.StatusBadRequest, ErrEmailExist)
		case repository.ErrDuplicateUsername:
			return nil, CreateServiceError(http.StatusBadRequest, ErrUsernameExist)
		default:
			s.logger.Errorf("Cannot InsertOne for %+v, Reason: %v", *d, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return u, nil
}

func (s *UserService) Update(ctx context.Context, d *dto.UpdateUserFields, uid string) (*model.User, *ServiceError) {
	s.logger.Infof("PUT Update User %#v userID: %s", d, uid)
	u, err := s.GetOneById(ctx, uid)
	if err != nil {
		return nil, err
	}

	if v := d.Password; v != nil {
		hashed := s.HashPassword(*v)
		d.Password = &hashed
	}

	if err := s.userRepo.UpdateOne(ctx, d, u); err != nil {
		switch err.Error() {
		case repository.ErrDuplicateEmail:
			return nil, CreateServiceError(http.StatusBadRequest, ErrEmailExist)
		case repository.ErrDuplicateUsername:
			return nil, CreateServiceError(http.StatusBadRequest, ErrUsernameExist)
		default:
			s.logger.Errorf("Cannot InsertOne for %#v, Reason: %v", d, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return u, nil
}

func (s *UserService) HashPassword(p string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("Cannot HashPassword, Reason: %v", err)
	}

	return string(hashed)
}

func (s *UserService) GetProfile(ctx context.Context, username, userID string) (*conduit.ProfileResponse, *ServiceError) {
	u, err := s.GetOne(ctx, &repository.FindOneUserFilter{Username: username})
	if err != nil {
		return nil, err
	}

	following := s.IsFollowing(ctx, userID, u.ID)

	res := &conduit.ProfileResponse{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}

	return res, nil
}
