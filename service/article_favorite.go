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
	log.Infof("POST FavoriteArticle user:%q, slug:%q", username, slug)
	a, err := s.GetArticleBySlug(ctx, username, slug)
	if err != nil {
		return nil, err
	}
	if err := s.favoritesRepo.InsertOne(ctx, username, a.ID); err != nil {
		log.Warnln("Cannot FavoriteArticle reason:", err)
		return nil, conduit.GeneralError
	}
	a.Favorited = true
	a.FavoritesCount += 1
	return a, nil
}

func (s *ArticleService) UnfavoriteArticleBySlug(ctx context.Context, username, slug string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("DELETE UnfavoriteArticle user:%q, slug:%q", username, slug)
	a, err := s.GetArticleBySlug(ctx, username, slug)
	if err != nil {
		return nil, err
	}
	if err := s.favoritesRepo.Delete(ctx, username, a.ID); err != nil {
		log.Warnln("Cannot UnfavoriteArticle reason:", err)
		return nil, conduit.GeneralError
	}
	a.Favorited = false
	a.FavoritesCount -= 1
	return a, nil
}

func (s *ArticleService) IsArticleFavorited(ctx context.Context, username, articleID string) bool {
	log := logger.GetCtx(ctx)
	if username == "" {
		return false
	}
	ptr, err := s.favoritesRepo.FindOneByIDs(ctx, username, articleID)
	if err != nil && err != sql.ErrNoRows {
		log.Warnln("Error get favorites repo", err)
	}
	return ptr != nil && err == nil
}
