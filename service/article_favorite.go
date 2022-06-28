package service

import (
	"context"
	"database/sql"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/utils/logger"
)

func (s *ArticleService) FavoriteArticleBySlug(ctx context.Context, username, slug string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST FavoriteArticle user_id: %s, slug: %s", username, slug)

	a, sErr := s.GetArticleBySlug(ctx, username, slug)
	if sErr != nil {
		return nil, sErr
	}

	if err := s.favoritesRepo.InsertOne(ctx, username, a.ID); err != nil {
		log.Warnf("Cannot FavoriteArticle, Reason: %v", err)
		return nil, conduit.GeneralError
	}

	a.Favorited = true
	count, err := s.favoritesRepo.CountFavorites(ctx, a.ID)
	if err != nil {
		log.Warnf("Cannot CountFavorites, Reason: %v", err)
		return nil, conduit.GeneralError
	}
	a.FavoritesCount = count
	return a, nil
}

func (s *ArticleService) UnfavoriteArticleBySlug(ctx context.Context, username, slug string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("DELETE UnfavoriteArticle user_id: %s, slug: %s", username, slug)

	a, err := s.GetArticleBySlug(ctx, username, slug)
	if err != nil {
		return nil, err
	}
	if err := s.favoritesRepo.Delete(ctx, username, a.ID); err != nil {
		log.Warnf("Cannot UnfavoriteArticle, Reason: %v", err)
		return nil, conduit.GeneralError
	}
	a.Favorited = false
	return a, nil
}

func (s *ArticleService) IsArticleFavorited(ctx context.Context, username, articleID string) bool {
	log := logger.GetCtx(ctx)
	if username == "" {
		return false
	}

	ptr, err := s.favoritesRepo.FindOneByIDs(ctx, username, articleID)
	if err != nil && err != sql.ErrNoRows {
		log.Warnf("Error get favorites repo %v", err)
	}
	return ptr != nil && err == nil
}
