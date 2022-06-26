package service

import (
	"context"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/utils/logger"
)

func (s *UserService) FollowUser(ctx context.Context, followerID, username string) (*model.ProfileResponse, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST FollowUser folowerID: %s, username: %s", followerID, username)

	following, err := s.GetOne(ctx, &repository.FindOneUserFilter{Username: username})
	if err != nil {
		return nil, err
	}

	if followerID == following.ID {
		return nil, conduit.BuildError(http.StatusBadRequest, ErrSelfFollow)
	}

	if err := s.followRepo.InsertOne(ctx, followerID, following.ID); err != nil {
		switch err.Error() {
		case repository.ErrDuplicateFollowing:
			return nil, conduit.BuildError(http.StatusBadRequest, ErrAlreadyFollow)
		default:
			log.Errorf("Cannot FollowUser followerID:%s following.ID:%s, Reason: %v", followerID, following.ID, err)
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

func (s *UserService) UnfollowUser(ctx context.Context, followerID, username string) (*model.ProfileResponse, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST UnfollowUser folowerID: %s, username: %s", followerID, username)

	following, err := s.GetOne(ctx, &repository.FindOneUserFilter{Username: username})
	if err != nil {
		return nil, err
	}

	if followerID == following.ID {
		return nil, conduit.BuildError(http.StatusBadRequest, ErrSelfUnfollow)
	}

	if err := s.followRepo.DeleteOneIDs(ctx, followerID, following.ID); err != nil {
		log.Errorf("Cannot UnfollowUser followerID:%s following.ID:%s, Reason: %v", followerID, following.ID, err)
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

func (s *UserService) IsFollowing(ctx context.Context, followerID, followingID string) bool {
	ptr, err := s.followRepo.FindOneByIDs(ctx, followerID, followingID)
	return ptr != nil && err == nil
}
