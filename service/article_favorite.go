package service

import (
	"context"
	"database/sql"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/utils/logger"
)

func (s *ArticleService) FavoriteArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST FavoriteArticle user_id: %s, slug: %s", userID, slug)

	a, sErr := s.GetArticleBySlug(ctx, userID, slug)
	if sErr != nil {
		return nil, sErr
	}

	if err := s.favoritesRepo.InsertOne(ctx, userID, a.ID); err != nil {
		log.Errorf("Cannot FavoriteArticle, Reason: %v", err)
		return nil, conduit.GeneralError
	}

	a.Favorited = true
	count, err := s.favoritesRepo.CountFavorites(ctx, a.ID)
	if err != nil {
		log.Errorf("Cannot CountFavorites, Reason: %v", err)
		return nil, conduit.GeneralError
	}

	a.FavoritesCount = count

	return a, nil
}

func (s *ArticleService) UnfavoriteArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("DELETE UnfavoriteArticle user_id: %s, slug: %s", userID, slug)

	a, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	if err := s.favoritesRepo.Delete(ctx, userID, a.ID); err != nil {
		log.Errorf("Cannot UnfavoriteArticle, Reason: %v", err)
		return nil, conduit.GeneralError
	}

	a.Favorited = false
	return a, nil
}

func (s *ArticleService) IsArticleFavorited(ctx context.Context, userID, articleID string) bool {
	log := logger.GetCtx(ctx)
	if userID == "" {
		return false
	}

	ptr, err := s.favoritesRepo.FindOneByIDs(ctx, userID, articleID)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("Error get favorites repo %v", err)
	}
	return ptr != nil && err == nil
}
