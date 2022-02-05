package service

import (
	"context"
	"net/http"

	"github.com/ashalfarhan/realworld/db/model"
)

func (s *ArticleService) FavoriteArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *ServiceError) {
	a, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return nil, err
	}

	if err := s.favoritesRepo.InsertOne(ctx, userID, a.ID); err != nil {
		s.logger.Printf("Cannot FavoriteArticle, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.Favorited = true
	return a, nil
}

func (s *ArticleService) UnfavoriteArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *ServiceError) {
	a, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	if err := s.favoritesRepo.Delete(ctx, userID, a.ID); err != nil {
		s.logger.Printf("Cannot UnfavoriteArticle, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.Favorited = false
	return a, nil
}

func (s *ArticleService) IsArticleFavorited(ctx context.Context, userID, articleID string) bool {
	ptr, err := s.favoritesRepo.FindOneByIDs(ctx, userID, articleID)
	return ptr != nil && err == nil
}
