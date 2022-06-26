package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/utils/logger"
)

func (s *ArticleService) CreateComment(ctx context.Context, d *model.CreateCommentDto) (*model.Comment, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST CreateComment %#v", d)

	ar, sErr := s.GetArticleBySlug(ctx, d.AuthorID, d.ArticleSlug)
	if sErr != nil {
		return nil, sErr
	}
	c := &model.Comment{
		Body:      d.Comment.Body,
		AuthorID:  d.AuthorID,
		ArticleID: ar.ID,
	}

	if err := s.commentRepo.InsertOne(ctx, c); err != nil {
		log.Printf("Cannot InsertOne::CommentRepo Args: %+v, Reason: %v", c, err)
		return nil, conduit.GeneralError
	}

	u, err := s.userRepo.FindOneByID(ctx, c.AuthorID)
	if err != nil {
		log.Printf("Cannot FindOneByID User Repo for %s, Reason: %v", c.AuthorID, err)
		return nil, conduit.GeneralError
	}

	c.Author = u
	return c, nil
}

func (s *ArticleService) GetComments(ctx context.Context, slug, userID string) ([]*model.Comment, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	ar, sErr := s.GetArticleBySlug(ctx, "", slug)
	if sErr != nil {
		return nil, sErr
	}

	comments, err := s.commentRepo.FindByArticleID(ctx, ar.ID)
	if err != nil {
		log.Printf("Cannot FindByArticleID::CommentRepo, Reason: %v", err)
		return nil, conduit.GeneralError
	}

	for _, c := range comments {
		u, err := s.userRepo.FindOneByID(ctx, c.AuthorID)
		if err != nil {
			log.Printf("Cannot FindOneByID User Repo for %s, Reason: %v", c.AuthorID, err)
			return nil, conduit.GeneralError
		}
		c.Author = u
	}

	return comments, nil
}

func (s *ArticleService) GetOneComment(ctx context.Context, commentID string) (*model.Comment, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	comm, err := s.commentRepo.FindOneByID(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, conduit.BuildError(http.StatusNotFound, ErrNoCommentFound)
		}

		log.Printf("Cannot FindOneByID::CommentRepo for %v, Reason: %v", commentID, err)
		return nil, conduit.GeneralError
	}

	return comm, nil
}

func (s *ArticleService) DeleteCommentByID(ctx context.Context, commentID, userID string) *model.ConduitError {
	log := logger.GetCtx(ctx)
	comm, err := s.GetOneComment(ctx, commentID)
	if err != nil {
		return err
	}

	if comm.AuthorID != userID {
		return conduit.BuildError(http.StatusForbidden, ErrNotAllowedDeleteComment)
	}

	if err := s.commentRepo.DeleteByID(ctx, commentID); err != nil {
		log.Printf("Cannot DeleteByID::CommentRepo, Reason: %v", err)
		return conduit.GeneralError
	}

	return nil
}
