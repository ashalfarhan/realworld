package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ashalfarhan/realworld/db/model"
)

func (s *ArticleService) FavoriteArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *ServiceError) {
	s.logger.Infof("POST FavoriteArticle user_id: %s, slug: %s", userID, slug)

	a, sErr := s.GetArticleBySlug(ctx, userID, slug)
	if sErr != nil {
		return nil, sErr
	}

	if err := s.favoritesRepo.InsertOne(ctx, userID, a.ID); err != nil {
		s.logger.Errorf("Cannot FavoriteArticle, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.Favorited = true
	count, err := s.favoritesRepo.CountFavorites(ctx, a.ID)
	if err != nil {
		s.logger.Errorf("Cannot CountFavorites, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.FavoritesCount = count

	return a, nil
}

func (s *ArticleService) UnfavoriteArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *ServiceError) {
	s.logger.Infof("DELETE UnfavoriteArticle user_id: %s, slug: %s", userID, slug)

	a, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	if err := s.favoritesRepo.Delete(ctx, userID, a.ID); err != nil {
		s.logger.Errorf("Cannot UnfavoriteArticle, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.Favorited = false
	return a, nil
}

func (s *ArticleService) IsArticleFavorited(ctx context.Context, userID, articleID string) bool {
	if userID == "" {
		return false
	}

	ptr, err := s.favoritesRepo.FindOneByIDs(ctx, userID, articleID)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("Error get favorites repo %v", err)
	}
	return ptr != nil && err == nil
}
