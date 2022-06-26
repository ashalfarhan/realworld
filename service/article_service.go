package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ashalfarhan/realworld/cache/store"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/gosimple/slug"
	"github.com/matoous/go-nanoid/v2"
)

const (
	defaultSlugId = 8
)

type ArticleService struct {
	articleRepo   repository.ArticleRepository
	userRepo      repository.UserRepository
	tagsRepo      repository.ArticleTagsRepository
	favoritesRepo repository.ArticleFavoritesRepository
	commentRepo   repository.CommentRepository
	articleCache  store.ArticleStore
}

func NewArticleService(repo *repository.Repository, store *store.CacheStore) *ArticleService {
	return &ArticleService{
		repo.ArticleRepo,
		repo.UserRepo,
		repo.ArticleTagsRepo,
		repo.ArticleFavoritesRepo,
		repo.CommentRepo,
		store.ArticleStore,
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, d *model.CreateArticleFields, authorID string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST CreateArticle %#v, author_id: %s", d, authorID)

	d.Slug = s.CreateSlug(d.Title)
	a, err := s.articleRepo.InsertOne(ctx, d, authorID)
	if err != nil {
		log.Printf("Cannot InsertOne to ArticleRepo for %+v, Reason: %v", a, err)
		return nil, conduit.GeneralError
	}

	if tgs := len(d.TagList); tgs > 0 {
		tags := make([]repository.InsertArticleTagsArgs, tgs)
		for i, tag := range d.TagList {
			tags[i] = repository.InsertArticleTagsArgs{ArticleID: a.ID, TagName: tag}
			a.TagList = append(a.TagList, tag)
		}
		if err := s.tagsRepo.InsertBulk(ctx, tags); err != nil {
			log.Errorf("Cannot InsertBulk::ArticleTags Repo for %v, Reason: %v", tags, err)
			return nil, conduit.GeneralError
		}
	}

	u, err := s.userRepo.FindOneByID(ctx, a.AuthorID)
	if err != nil {
		log.Errorf("Cannot FindOneByID User Repo for %s, Reason: %v", a.AuthorID, err)
		return nil, conduit.GeneralError
	}

	a.Author = u
	return a, nil
}

func (s *ArticleService) GetArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)

	if cached := s.articleCache.FindOneBySlug(ctx, slug, userID); cached != nil {
		return cached, nil
	}

	ar, err := s.articleRepo.FindOneBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, conduit.BuildError(http.StatusNotFound, ErrNoArticleFound)
		}
		log.Printf("Failed to get article by slug: %v", err)
		return nil, conduit.GeneralError
	}

	if err := s.PopulateArticleField(ctx, ar, userID); err != nil {
		return nil, err
	}

	s.articleCache.SaveBySlug(ctx, slug, userID, ar)
	return ar, nil
}

func (s *ArticleService) GetArticles(ctx context.Context, args *repository.FindArticlesArgs) (model.Articles, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	var articles model.Articles
	var err error

	// if cached := s.articleCache.Find(ctx, args); cached != nil {
	// 	return cached, nil
	// }

	articles, err = s.articleRepo.Find(ctx, args)
	if err != nil {
		log.Errorf("Cannot Find::ArticleRepo args: %+v, Reason: %v", args, err)
		return nil, conduit.GeneralError
	}

	for _, a := range articles {
		if err := s.PopulateArticleField(ctx, a, args.UserID); err != nil {
			return nil, err
		}
	}

	// s.caching.Set(ctx, cacheKey, articles)
	return articles, nil
}

func (s *ArticleService) GetArticlesFeed(ctx context.Context, args *repository.FindArticlesArgs) (model.Articles, *model.ConduitError) {
	log := logger.GetCtx(ctx)

	// cacheKey := fmt.Sprintf("articles-feed|user_id:%s|limit:%d|offset:%d", args.UserID, args.Limit, args.Offset)
	// if ok := s.caching.Get(ctx, cacheKey, &articles); ok {
	// 	log.Infof("Response with cache %s", cacheKey)
	// 	return articles, nil
	// }

	articles, err := s.articleRepo.FindByFollowed(ctx, args)
	if err != nil {
		log.Errorf("Cannot FindByFollowed::ArticleRepo args: %+v, Reason: %v", args, err)
		return nil, conduit.GeneralError
	}

	for _, a := range articles {
		if err := s.PopulateArticleField(ctx, a, args.UserID); err != nil {
			return nil, err
		}
	}

	// s.caching.Set(ctx, cacheKey, articles)
	return articles, nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, slug, userID string) *model.ConduitError {
	log := logger.GetCtx(ctx)

	a, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return err
	}
	if a.AuthorID != userID {
		log.Warnf("Forbidden delete article author_id:%q, user_id: %q", a.AuthorID, userID)
		return conduit.BuildError(http.StatusForbidden, ErrNotAllowedDeleteArticle)
	}
	if err := s.articleRepo.DeleteBySlug(ctx, slug); err != nil {
		log.Printf("Failed to delete article by slug:%q, Reason: %v", slug, err)
		return conduit.GeneralError
	}
	return nil
}

func (s *ArticleService) UpdateArticleBySlug(ctx context.Context, userID, slug string, d *model.UpdateArticleFields) (*model.Article, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("UpdateArticleBySlug user_id: %s, slug: %s, %#v", userID, slug, d)

	ar, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return nil, err
	}

	if ar.AuthorID != userID {
		return nil, conduit.BuildError(http.StatusForbidden, ErrNotAllowedUpdateArticle)
	}

	// Updating title will update the slug
	if v := d.Title; v != nil {
		newSlug := s.CreateSlug(*v)
		d.Slug = &newSlug
	}

	if err := s.articleRepo.UpdateOneBySlug(ctx, d, ar); err != nil {
		log.Errorf("Cannot UpdateOneBySlug, slug: %s, payload: %+v, Reason: %v", slug, d, err)
		return nil, conduit.GeneralError
	}

	return ar, nil
}

func (s *ArticleService) PopulateArticleField(ctx context.Context, a *model.Article, userID string) *model.ConduitError {
	log := logger.GetCtx(ctx)
	tags, err := s.tagsRepo.FindArticleTagsByID(ctx, a.ID)
	if err != nil {
		log.Errorf("Cannot FindArticleTagsByID::ArticleTagsRepo for %s, Reason: %v", a.ID, err)
		return conduit.GeneralError
	}

	a.TagList = tags

	a.Author, err = s.userRepo.FindOneByID(ctx, a.AuthorID)
	if err != nil {
		log.Errorf("Cannot FindOneByID::UserRepo for %s, Reason: %v", a.AuthorID, err)
		return conduit.GeneralError
	}

	a.Favorited = s.IsArticleFavorited(ctx, userID, a.ID)

	a.FavoritesCount, err = s.favoritesRepo.CountFavorites(ctx, a.ID)
	if err != nil {
		log.Errorf("Cannot CountFavorites FavoritesRepo for %s, Reason: %v", a.ID, err)
		return conduit.GeneralError
	}

	return nil
}

// TODO: Can be moved to utils
func (s *ArticleService) CreateSlug(title string) string {
	id, _ := gonanoid.New(defaultSlugId)
	slug := slug.Make(title)
	return fmt.Sprintf("%s-%s", slug, id)
}
