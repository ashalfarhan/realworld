package service

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   repository.UserRepository
	followRepo repository.FollowingRepository
	logger     *log.Logger
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		repo.UserRepo,
		repo.FollowRepo,
		conduit.NewLogger("user-service"),
	}
}

func (s *UserService) GetOneById(ctx context.Context, id string) (*model.User, *ServiceError) {
	u := &model.User{}
	if err := s.userRepo.FindOneByID(ctx, id, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, ErrNoUserFound)
		}

		s.logger.Printf("Cannot FindOneById for %s, Reason: %v", id, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return u, nil
}

type GetOneArgs struct {
	Email    string
	Username string
	UserID   string
}

func (s *UserService) GetOne(ctx context.Context, d *GetOneArgs) (*model.User, *ServiceError) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}

	if err := s.userRepo.FindOne(ctx, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, ErrNoUserFound)
		}

		s.logger.Printf("Cannot FindOne for %+v, Reason: %v", d, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return u, nil
}

type RegisterArgs struct {
	Email    string
	Username string
	Password string
}

func (s *UserService) Register(ctx context.Context, d *RegisterArgs) (*model.User, *ServiceError) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}

	u.Password = s.HashPassword(d.Password)

	if err := s.userRepo.InsertOne(ctx, u); err != nil {
		switch err.Error() {
		case repository.ErrDuplicateEmail:
			return nil, CreateServiceError(http.StatusBadRequest, ErrEmailExist)
		case repository.ErrDuplicateUsername:
			return nil, CreateServiceError(http.StatusBadRequest, ErrUsernameExist)
		default:
			s.logger.Printf("Cannot InsertOne for %+v, Reason: %v", d, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}

	}

	return u, nil
}

func (s *UserService) Update(ctx context.Context, d *dto.UpdateUserDto, uid string) (*model.User, *ServiceError) {
	u, err := s.GetOneById(ctx, uid)
	if err != nil {
		return nil, err
	}

	args := &repository.UpdateUserValues{
		ID:    uid,
		Bio:   d.User.Bio,
		Image: d.User.Image,
	}

	if len(d.User.Email) != 0 {
		args.Email = &d.User.Email
		u.Email = d.User.Email
	}
	if len(d.User.Username) != 0 {
		args.Username = &d.User.Username
		u.Username = d.User.Username
	}
	if len(d.User.Password) != 0 {
		hashed := s.HashPassword(d.User.Password)
		args.Password = &hashed
		u.Password = hashed
	}

	if d.User.Bio.Set {
		args.Bio = d.User.Bio
		u.Bio = d.User.Bio
	}

	if d.User.Image.Set {
		args.Image = d.User.Image
		u.Image = d.User.Image
	}

	if err := s.userRepo.UpdateOne(ctx, args); err != nil {
		switch err.Error() {
		case repository.ErrDuplicateEmail:
			return nil, CreateServiceError(http.StatusBadRequest, ErrEmailExist)
		case repository.ErrDuplicateUsername:
			return nil, CreateServiceError(http.StatusBadRequest, ErrUsernameExist)
		default:
			s.logger.Printf("Cannot UpdateOne payload:%+v args:%+v, Reason: %v", d, args, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return u, nil
}

func (s *UserService) HashPassword(p string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Printf("Cannot HashPassword, Reason: %v", err)
	}

	return string(hashed)
}
