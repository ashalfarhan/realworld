package service

import (
	"context"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/utils/logger"
)

func (s *UserService) FollowUser(ctx context.Context, followUsername, username string) (*model.ProfileRs, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST FollowUser followUsername:%q, user:%q", followUsername, username)
	following, err := s.GetOne(ctx, &model.FindUserArg{Username: username})
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
			log.Warnln("Cannot insert to follow repo reason:", err)
			return nil, conduit.GeneralError
		}
	}

	res := &model.ProfileRs{
		Username:  following.Username,
		Bio:       following.Bio,
		Image:     following.Image,
		Following: true,
	}
	return res, nil
}

func (s *UserService) UnfollowUser(ctx context.Context, followUsername, username string) (*model.ProfileRs, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST UnfollowUser followUsername:%q, user:%q", followUsername, username)
	following, err := s.GetOne(ctx, &model.FindUserArg{Username: username})
	if err != nil {
		return nil, err
	}

	if followUsername == following.Username {
		return nil, conduit.BuildError(http.StatusBadRequest, ErrSelfUnfollow)
	}

	if err := s.followRepo.DeleteOneIDs(ctx, followUsername, following.Username); err != nil {
		log.Warnln("Cannot delete to follow repo reason:", followUsername, following.ID, err)
		return nil, conduit.GeneralError
	}

	res := &model.ProfileRs{
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
