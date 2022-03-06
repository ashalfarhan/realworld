package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/cache"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/gosimple/slug"
	"github.com/matoous/go-nanoid/v2"
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
	logger        *log.Logger
	caching       *CacheService
}

func NewArticleService(repo *repository.Repository) *ArticleService {
	return &ArticleService{
		repo.ArticleRepo,
		repo.UserRepo,
		repo.ArticleTagsRepo,
		repo.ArticleFavoritesRepo,
		repo.CommentRepo,
		conduit.NewLogger("article-service"),
		NewCacheService(cache.Ca),
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, d *dto.CreateArticleDto, authorID string) (*model.Article, *ServiceError) {
	a := &model.Article{
		Title:       d.Article.Title,
		Description: d.Article.Description,
		Body:        d.Article.Body,
		AuthorID:    authorID,
	}

	a.Slug = s.CreateSlug(a.Title)

	if err := s.articleRepo.InsertOne(ctx, a); err != nil {
		s.logger.Printf("Cannot InsertOne to ArticleRepo for %+v, Reason: %v", a, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	if len(d.Article.TagList) > 0 {
		for _, tag := range d.Article.TagList {
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

func (s *ArticleService) GetArticleBySlug(ctx context.Context, userID string, slug string) (*model.Article, *ServiceError) {
	a := &model.Article{
		Slug:   slug,
		Author: &model.User{},
	}

	cacheKey := fmt.Sprintf("article-%s-%s", slug, userID)
	if ok := s.caching.Get(ctx, cacheKey, a); ok {
		return a, nil
	}

	if err := s.articleRepo.FindOneBySlug(ctx, a); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, ErrNoArticleFound)
		}

		s.logger.Printf("Cannot FindOneBySlug Article Repo for %+v, Reason: %v", a, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	s.PopulateArticleField(ctx, a, userID)
	s.caching.Set(ctx, cacheKey, a)
	return a, nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, slug, userID string) *ServiceError {
	a, err := s.GetArticleBySlug(ctx, userID, slug)
	if err != nil {
		return err
	}

	if a.AuthorID != userID {
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
		s.logger.Printf("Cannot create nanoid for %s, Reason: %v", title, err)
	}

	slug := slug.Make(title)

	return fmt.Sprintf("%s-%s", slug, id)
}

func (s *ArticleService) UpdateArticleBySlug(ctx context.Context, userID, slug string, d *dto.UpdateArticleDto) (*model.Article, *ServiceError) {
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

func (s *ArticleService) GetArticles(ctx context.Context, args *repository.FindArticlesArgs) (model.Articles, *ServiceError) {
	var articles model.Articles
	var err error
	cacheKey := fmt.Sprintf("articles-list-%s-%d-%d", args.UserID, args.Limit, args.Offset)

	if ok := s.caching.Get(ctx, cacheKey, &articles); ok {
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

func (s *ArticleService) GetArticlesFeed(ctx context.Context, args *repository.FindArticlesArgs) (model.Articles, *ServiceError) {
	var articles model.Articles
	var err error
	cacheKey := fmt.Sprintf("articles-feed-%s-%d-%d", args.UserID, args.Limit, args.Offset)

	if ok := s.caching.Get(ctx, cacheKey, &articles); ok {
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

func (s *ArticleService) PopulateArticleField(ctx context.Context, a *model.Article, userID string) *ServiceError {
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
