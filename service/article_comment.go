package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
)

func (s *ArticleService) CreateComment(ctx context.Context, d *dto.CreateCommentDto) (*model.Comment, *ServiceError) {
	ar, err := s.GetArticleBySlug(ctx, d.AuthorID, d.ArticleSlug)
	if err != nil {
		return nil, err
	}
	c := &model.Comment{
		Body:      d.Comment.Body,
		AuthorID:  d.AuthorID,
		ArticleID: ar.ID,
	}

	if err := s.commentRepo.InsertOne(ctx, c); err != nil {
		s.logger.Printf("Cannot InsertOne::CommentRepo Args: %+v, Reason: %v", c, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	c.Author = &model.User{}
	if err := s.userRepo.FindOneByID(ctx, c.AuthorID, c.Author); err != nil {
		s.logger.Printf("Cannot FindOneByID User Repo for %s, Reason: %v", c.AuthorID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return c, nil
}

func (s *ArticleService) GetComments(ctx context.Context, args *repository.FindCommentsByArticleIDArgs, slug string) ([]*model.Comment, *ServiceError) {
	ar, sErr := s.GetArticleBySlug(ctx, "", slug)
	if sErr != nil {
		return nil, sErr
	}

	args.ArticleID = ar.ID
	comments, err := s.commentRepo.FindByArticleID(ctx, args)
	if err != nil {
		s.logger.Printf("Cannot FindByArticleID::CommentRepo, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	for _, c := range comments {
		c.Author = &model.User{}
		if err := s.userRepo.FindOneByID(ctx, c.AuthorID, c.Author); err != nil {
			s.logger.Printf("Cannot FindOneByID User Repo for %s, Reason: %v", c.AuthorID, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return comments, nil
}

func (s *ArticleService) GetOneComment(ctx context.Context, commentID string) (*model.Comment, *ServiceError) {
	comm, err := s.commentRepo.FindOneByID(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, ErrNoCommentFound)
		}

		s.logger.Printf("Cannot FindOneByID::CommentRepo for %v, Reason: %v", commentID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return comm, nil
}

func (s *ArticleService) DeleteCommentByID(ctx context.Context, commentID, userID string) *ServiceError {
	comm, err := s.GetOneComment(ctx, commentID)
	if err != nil {
		return err
	}

	if comm.AuthorID != userID {
		return CreateServiceError(http.StatusForbidden, ErrNotAllowedDeleteComment)
	}

	if err := s.commentRepo.DeleteByID(ctx, commentID); err != nil {
		s.logger.Printf("Cannot DeleteByID::CommentRepo, Reason: %v", err)
		return CreateServiceError(http.StatusInternalServerError, nil)
	}

	return nil
}
