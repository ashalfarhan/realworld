package service

import (
	"context"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/utils/logger"
)

func (s *UserService) FollowUser(ctx context.Context, followUsername, username string) (*model.ProfileResponse, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST FollowUser folowerID: %s, username: %s", followUsername, username)

	following, err := s.GetOne(ctx, &repository.FindOneUserFilter{Username: username})
	if err != nil {
		return nil, err
	}

	if followUsername == following.Username {
		return nil, conduit.BuildError(http.StatusBadRequest, ErrSelfFollow)
	}

	if err := s.followRepo.InsertOne(ctx, followUsername, following.Username); err != nil {
		switch err.Error() {
		case repository.ErrDuplicateFollowing:
			return nil, conduit.BuildError(http.StatusBadRequest, ErrAlreadyFollow)
		default:
			log.Warnf("Cannot FollowUser followerID:%s following.ID:%s, Reason: %v", followUsername, following.ID, err)
			return nil, conduit.GeneralError
		}
	}

	res := &model.ProfileResponse{
		Username:  following.Username,
		Bio:       following.Bio,
		Image:     following.Image,
		Following: true,
	}
	return res, nil
}

func (s *UserService) UnfollowUser(ctx context.Context, followUsername, username string) (*model.ProfileResponse, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST UnfollowUser folowerID: %s, username: %s", followUsername, username)

	following, err := s.GetOne(ctx, &repository.FindOneUserFilter{Username: username})
	if err != nil {
		return nil, err
	}

	if followUsername == following.Username {
		return nil, conduit.BuildError(http.StatusBadRequest, ErrSelfUnfollow)
	}

	if err := s.followRepo.DeleteOneIDs(ctx, followUsername, following.Username); err != nil {
		log.Warnf("Cannot UnfollowUser followerID:%s following.ID:%s, Reason: %v", followUsername, following.ID, err)
		return nil, conduit.GeneralError
	}

	res := &model.ProfileResponse{
		Username:  following.Username,
		Bio:       following.Bio,
		Image:     following.Image,
		Following: false,
	}
	return res, nil
}

func (s *UserService) IsFollowing(ctx context.Context, followerUsername, followingUsername string) bool {
	ptr, err := s.followRepo.FindOneByIDs(ctx, followerUsername, followingUsername)
	return ptr != nil && err == nil
}
