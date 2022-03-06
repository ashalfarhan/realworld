package service

import (
	"context"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/repository"
)

func (s *UserService) FollowUser(ctx context.Context, followerID, username string) (*conduit.ProfileResponse, *ServiceError) {
	following, err := s.GetOne(ctx, &repository.FindOneUserFilter{Username: username})
	if err != nil {
		return nil, err
	}

	if followerID == following.ID {
		return nil, CreateServiceError(http.StatusBadRequest, ErrSelfFollow)
	}

	if err := s.followRepo.InsertOne(ctx, followerID, following.ID); err != nil {
		switch err.Error() {
		case repository.ErrDuplicateFollowing:
			return nil, CreateServiceError(http.StatusBadRequest, ErrAlreadyFollow)
		default:
			s.logger.Printf("Cannot FollowUser followerID:%s following.ID:%s, Reason: %v", followerID, following.ID, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	res := &conduit.ProfileResponse{
		Username:  following.Username,
		Bio:       following.Bio,
		Image:     following.Image,
		Following: true,
	}

	return res, nil
}

func (s *UserService) UnfollowUser(ctx context.Context, followerID, username string) (*conduit.ProfileResponse, *ServiceError) {
	following, err := s.GetOne(ctx, &repository.FindOneUserFilter{Username: username})
	if err != nil {
		return nil, err
	}

	if followerID == following.ID {
		return nil, CreateServiceError(http.StatusBadRequest, ErrSelfUnfollow)
	}

	if err := s.followRepo.DeleteOneIDs(ctx, followerID, following.ID); err != nil {
		s.logger.Printf("Cannot UnfollowUser followerID:%s following.ID:%s, Reason: %v", followerID, following.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, err)
	}

	res := &conduit.ProfileResponse{
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
