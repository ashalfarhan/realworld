package service

import (
	"database/sql"
	"errors"

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
	return &UserService{repo.UR, repo.FR}
}

func (s *UserService) GetOneById(id string) (*model.User, error) {
	u := &model.User{}
	if err := s.userRepo.FindOneById(id, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no user found")
		}
		return nil, conduit.ErrInternal
	}

	return u, nil
}

func (s *UserService) GetOne(d *dto.LoginUserDto) (*model.User, error) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}
	if err := s.userRepo.FindOne(u); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no user found")
		}
		return nil, conduit.ErrInternal
	}

	return u, nil
}

func (s *UserService) CreateOne(d *dto.RegisterUserDto) (*model.User, error) {
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
	}
	if err := s.userRepo.FindOne(u); err == nil {
		return nil, errors.New("username or email is already used")
	}

	if err := u.HashPassword(d.Password); err != nil {
		return nil, conduit.ErrInternal
	}

	if err := s.userRepo.InsertOne(u); err != nil {
		return nil, conduit.ErrInternal
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

func (s *UserService) FollowUser(followerID, username string) error {
	following, err := s.GetOne(&dto.LoginUserDto{Username: username})
	if err != nil {
		return err
	}

	if followerID == following.ID {
		return errors.New("you cannot follow your self")
	}

	return s.followRepo.Follow(followerID, following.ID)
}

func (s *UserService) UnfollowUser(followerID, username string) error {
	following, err := s.GetOne(&dto.LoginUserDto{Username: username})
	if err != nil {
		return err
	}

	if followerID == following.ID {
		return errors.New("you cannot unfollow your self")
	}

	return s.followRepo.Unfollow(followerID, following.ID)
}

func (s *UserService) IsFollowing(followerID, followingID string) bool {
	return s.followRepo.IsFollowing(followerID, followingID)
}
