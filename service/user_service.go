package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/utils/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   repository.UserRepository
	followRepo repository.FollowingRepository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		userRepo:   repo.UserRepo,
		followRepo: repo.FollowRepo,
	}
}

func (s *UserService) GetOneByUsername(ctx context.Context, username string) (*model.User, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	u, err := s.userRepo.FindOneByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, conduit.BuildError(http.StatusNotFound, ErrNoUserFound)
		}
		log.Errorf("Cannot FindOneById for %s, reason: %v", username, err)
		return nil, conduit.GeneralError
	}
	return u, nil
}

func (s *UserService) GetOne(ctx context.Context, d *repository.FindOneUserFilter) (*model.User, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	u, err := s.userRepo.FindOne(ctx, d)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, conduit.BuildError(http.StatusNotFound, ErrNoUserFound)
		}
		log.Errorf("Cannot find one in user repo for filter:%+v reason: %v", d, err)
		return nil, conduit.GeneralError
	}
	return u, nil
}

func (s *UserService) Insert(ctx context.Context, d *model.RegisterUserFields) (*model.User, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	d.Password = s.HashPassword(d.Password)
	u, err := s.userRepo.InsertOne(ctx, d)
	if err != nil {
		switch err.Error() {
		case repository.ErrDuplicateEmail:
			return nil, conduit.BuildError(http.StatusBadRequest, ErrEmailExist)
		case repository.ErrDuplicateUsername:
			return nil, conduit.BuildError(http.StatusBadRequest, ErrUsernameExist)
		default:
			log.Errorf("Cannot insert to user repo reason: %v", err)
			return nil, conduit.GeneralError
		}
	}
	return u, nil
}

func (s *UserService) Update(ctx context.Context, d *model.UpdateUserFields, username string) (*model.User, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	u, err := s.GetOneByUsername(ctx, username)
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
			return nil, conduit.BuildError(http.StatusBadRequest, ErrEmailExist)
		case repository.ErrDuplicateUsername:
			return nil, conduit.BuildError(http.StatusBadRequest, ErrUsernameExist)
		default:
			log.Errorf("Cannot InsertOne for %+v, Reason: %v", d, err)
			return nil, conduit.GeneralError
		}
	}
	return u, nil
}

func (s *UserService) HashPassword(p string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hashed)
}

func (s *UserService) GetProfile(ctx context.Context, username, userID string) (*model.ProfileResponse, *model.ConduitError) {
	u, err := s.GetOne(ctx, &repository.FindOneUserFilter{Username: username})
	if err != nil {
		return nil, err
	}
	following := s.IsFollowing(ctx, userID, u.ID)
	res := &model.ProfileResponse{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}
	return res, nil
}
