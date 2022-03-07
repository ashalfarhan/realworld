package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/cache"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/gosimple/slug"
	"github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
)

const (
	defaultSlugId = 8
)

type ArticleService struct {
	articleRepo   repository.ArticleRepository
	userRepo      repository.UserRepository
	tagsRepo      *repository.ArticleTagsRepository
	favoritesRepo repository.ArticleFavoritesRepository
	commentRepo   repository.CommentRepository
	logger        *logrus.Entry
	caching       *CacheService
}

func NewArticleService(repo *repository.Repository) *ArticleService {
	return &ArticleService{
		repo.ArticleRepo,
		repo.UserRepo,
		repo.ArticleTagsRepo,
		repo.ArticleFavoritesRepo,
		repo.CommentRepo,
		conduit.NewLogger("service", "ArticleService"),
		NewCacheService(cache.Ca),
	}
}

func (s *ArticleService) CreateArticle(
	ctx context.Context, d *dto.CreateArticleFields, authorID string,
) (*model.Article, *ServiceError) {
	s.logger.Infof("CreateArticle %#v, author_id: %s", d, authorID)

	d.Slug = s.CreateSlug(d.Title)
	a, err := s.articleRepo.InsertOne(ctx, d, authorID)
	if err != nil {
		s.logger.Printf("Cannot InsertOne to ArticleRepo for %+v, Reason: %v", a, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	if len(d.TagList) > 0 {
		for _, tag := range d.TagList {
			if err := s.tagsRepo.InsertOne(ctx, a.ID, tag); err != nil {
				s.logger.Printf("Cannot InsertOne ArticleTags Repo for %s, Reason: %v", tag, err)
				return nil, CreateServiceError(http.StatusInternalServerError, nil)
			}
			a.TagList = append(a.TagList, tag)
		}
	}

	u, err := s.userRepo.FindOneByID(ctx, a.AuthorID)
	if err != nil {
		s.logger.Printf("Cannot FindOneByID User Repo for %s, Reason: %v", a.AuthorID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.Author = u
	return a, nil
}

func (s *ArticleService) GetArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *ServiceError) {
	a := &model.Article{
		Slug:   slug,
		Author: new(model.User),
	}

	cacheKey := fmt.Sprintf("article|slug:%s|user_id:%s", slug, userID)
	if ok := s.caching.Get(ctx, cacheKey, a); ok {
		s.logger.Infof("Response with cache %s", cacheKey)
		return a, nil
	}

	ar, err := s.articleRepo.FindOneBySlug(ctx, a.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, ErrNoArticleFound)
		}

		s.logger.Printf("Cannot FindOneBySlug Article Repo for %+v, Reason: %v", a, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a = ar
	s.PopulateArticleField(ctx, a, userID)
	s.caching.Set(ctx, cacheKey, a)
	return a, nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, slug, userID string) *ServiceError {
	s.logger.Infof("DeleteArticle slug: %s user_id: %s", slug, userID)
	a, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return err
	}

	if a.AuthorID != userID {
		s.logger.Warnf("Forbidden delete article author_id: %s, user_id: %s", a.AuthorID, userID)
		return CreateServiceError(http.StatusForbidden, ErrNotAllowedDeleteArticle)
	}

	if err := s.articleRepo.DeleteBySlug(ctx, slug); err != nil {
		s.logger.Printf("Cannot DeleteArticleBySlug Article Repo for %s, Reason: %v", slug, err)
		return CreateServiceError(http.StatusInternalServerError, nil)
	}

	return nil
}

func (s *ArticleService) CreateSlug(title string) string {
	id, err := gonanoid.New(defaultSlugId)
	if err != nil {
		s.logger.Errorf("Cannot create nanoid for %s, Reason: %v", title, err)
	}

	slug := slug.Make(title)
	return fmt.Sprintf("%s-%s", slug, id)
}

func (s *ArticleService) UpdateArticleBySlug(
	ctx context.Context, userID, slug string, d *dto.UpdateArticleDto,
) (*model.Article, *ServiceError) {
	s.logger.Infof("UpdateArticleBySlug user_id: %s, slug: %s, %#v", userID, slug, d)
	args := &repository.UpdateArticleValues{}
	ar, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return nil, err
	}

	if len(d.Article.Body) != 0 {
		args.Body = &d.Article.Body
		ar.Body = d.Article.Body
	}

	if len(d.Article.Description) != 0 {
		args.Description = &d.Article.Description
		ar.Description = d.Article.Description
	}

	if len(d.Article.Title) != 0 {
		args.Title = &d.Article.Title
		newSlug := s.CreateSlug(d.Article.Title)
		args.Slug = &newSlug

		ar.Title = d.Article.Title
		ar.Slug = newSlug
	}
	if ar.AuthorID != userID {
		return nil, CreateServiceError(http.StatusForbidden, ErrNotAllowedUpdateArticle)
	}

	if err := s.articleRepo.UpdateOneBySlug(ctx, slug, args, ar); err != nil {
		s.logger.Printf("Cannot UpdateOneBySlug, slug: %s, payload: %+v, Reason: %v", slug, d, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return ar, nil
}

func (s *ArticleService) GetArticles(
	ctx context.Context, args *repository.FindArticlesArgs,
) (model.Articles, *ServiceError) {
	var articles model.Articles
	var err error

	cacheKey := fmt.Sprintf("articles-list|user_id:%s|limit:%d|offset:%d", args.UserID, args.Limit, args.Offset)
	if ok := s.caching.Get(ctx, cacheKey, &articles); ok {
		s.logger.Infof("Response with cache %s", cacheKey)
		return articles, nil
	}

	articles, err = s.articleRepo.Find(ctx, args)
	if err != nil {
		s.logger.Printf("Cannot Find::ArticleRepo args: %+v, Reason: %v", args, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	for _, a := range articles {
		s.PopulateArticleField(ctx, a, args.UserID)
	}

	s.caching.Set(ctx, cacheKey, articles)
	return articles, nil
}

func (s *ArticleService) GetArticlesFeed(
	ctx context.Context, args *repository.FindArticlesArgs,
) (model.Articles, *ServiceError) {
	var articles model.Articles
	var err error

	cacheKey := fmt.Sprintf("articles-feed|user_id:%s|limit:%d|offset:%d", args.UserID, args.Limit, args.Offset)
	if ok := s.caching.Get(ctx, cacheKey, &articles); ok {
		s.logger.Infof("Response with cache %s", cacheKey)
		return articles, nil
	}

	articles, err = s.articleRepo.FindByFollowed(ctx, args)
	if err != nil {
		s.logger.Printf("Cannot FindByFollowed::ArticleRepo args: %+v, Reason: %v", args, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	for _, a := range articles {
		s.PopulateArticleField(ctx, a, args.UserID)
	}

	s.caching.Set(ctx, cacheKey, articles)
	return articles, nil
}

func (s *ArticleService) PopulateArticleField(
	ctx context.Context, a *model.Article, userID string,
) *ServiceError {
	tags, err := s.tagsRepo.FindArticleTagsByID(ctx, a.ID)
	if err != nil {
		s.logger.Printf("Cannot FindArticleTagsByID::ArticleTagsRepo for %s, Reason: %v", a.ID, err)
		return CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.TagList = tags

	u, err := s.userRepo.FindOneByID(ctx, a.AuthorID)
	if err != nil {
		s.logger.Printf("Cannot FindOneByID::UserRepo for %s, Reason: %v", a.AuthorID, err)
		return CreateServiceError(http.StatusInternalServerError, nil)
	}
	a.Author = u

	a.Favorited = s.IsArticleFavorited(ctx, userID, a.ID)

	a.FavoritesCount, err = s.favoritesRepo.CountFavorites(ctx, a.ID)
	if err != nil {
		s.logger.Printf("Cannot CountFavorites FavoritesRepo for %s, Reason: %v", a.ID, err)
		return CreateServiceError(http.StatusInternalServerError, nil)
	}

	return nil
}
