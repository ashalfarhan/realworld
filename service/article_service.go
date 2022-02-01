package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
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
	articleRepo   *repository.ArticleRepository
	userRepo      *repository.UserRepository
	tagsRepo      *repository.ArticleTagsRepository
	favoritesRepo *repository.ArticleFavoritesRepository
	logger        *log.Logger
}

func NewArticleService(repo *repository.Repository) *ArticleService {
	return &ArticleService{
		repo.ArticleRepo,
		repo.UserRepo,
		repo.ArticleTagsRepo,
		repo.ArticleFavoritesRepo,
		conduit.NewLogger("article-service"),
	}
}

func (s *ArticleService) Create(ctx context.Context, d *dto.CreateArticleDto, authorID string) (*model.Article, *ServiceError) {
	a := &model.Article{
		Title:       d.Title,
		Description: d.Description,
		Body:        d.Body,
		Author: &model.User{
			ID: authorID,
		},
	}

	a.Slug = s.CreateSlug(a.Title)

	if err := s.articleRepo.InsertOne(ctx, a); err != nil {
		s.logger.Printf("Cannot InsertOne to ArticleRepo for %#v, Reason: %v", a, err)
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

	if err := s.userRepo.FindOneByID(ctx, a.Author.ID, a.Author); err != nil {
		s.logger.Printf("Cannot FindOneByID User Repo for %s, Reason: %v", a.Author.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return a, nil
}

func (s *ArticleService) GetOneBySlug(ctx context.Context, slug string) (*model.Article, *ServiceError) {
	a := &model.Article{
		Slug:   slug,
		Author: &model.User{},
	}

	if err := s.articleRepo.FindOneBySlug(ctx, a); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, errors.New("no article found"))
		}
		s.logger.Printf("Cannot FindOneBySlug Article Repo for %#v, Reason: %v", a, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	tags, err := s.tagsRepo.GetArticleTagsByID(ctx, a.ID)
	if err != nil {
		s.logger.Printf("Cannot GetArticleTagsByID ArticleTags Repo for %s, Reason: %v", a.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.TagList = tags

	if err := s.userRepo.FindOneByID(ctx, a.Author.ID, a.Author); err != nil {
		s.logger.Printf("Cannot FindOneByID User Repo for %s, Reason: %v", a.Author.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return a, nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, slug string, userID string) *ServiceError {
	a, err := s.GetOneBySlug(ctx, slug)
	if err != nil {
		return err
	}

	if a.Author.ID != userID {
		return CreateServiceError(http.StatusForbidden, errors.New("you cannot delete this article"))
	}

	if err := s.articleRepo.DeleteBySlug(ctx, slug); err != nil {
		s.logger.Printf("Cannot DeleteArticleBySlug Article Repo for %s, Reason: %v", slug, err)
		return CreateServiceError(http.StatusInternalServerError, nil)
	}

	return nil
}

func (s *ArticleService) GetAllTags(ctx context.Context) ([]string, *ServiceError) {
	tags, err := s.tagsRepo.GetAllTags(ctx)
	if err != nil {
		s.logger.Printf("Cannot GetAllTags, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return tags, nil
}

func (s *ArticleService) CreateSlug(title string) string {
	id, err := gonanoid.New(defaultSlugId)
	if err != nil {
		s.logger.Printf("Cannot create nanoid for %s, Reason: %v", title, err)
	}

	slug := slug.Make(title)

	return fmt.Sprintf("%s-%s", slug, id)
}

func (s *ArticleService) UpdateOneBySlug(ctx context.Context, userID, slug string, d *dto.UpdateArticleDto) (*model.Article, *ServiceError) {
	args := &conduit.UpdateArticleArgs{}
	ar, err := s.GetOneBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if len(d.Body) != 0 {
		args.Body = &d.Body
		ar.Body = d.Body
	}

	if len(d.Description) != 0 {
		args.Description = &d.Description
		ar.Description = d.Description
	}

	if len(d.Title) != 0 {
		args.Title = &d.Title
		newSlug := s.CreateSlug(d.Title)
		args.Slug = &newSlug

		ar.Title = d.Title
		ar.Slug = newSlug
	}

	if ar.Author.ID != userID {
		return nil, CreateServiceError(http.StatusForbidden, errors.New("you cannot edit this article"))
	}

	if err := s.articleRepo.UpdateOneBySlug(ctx, slug, args, ar); err != nil {
		s.logger.Printf("Cannot UpdateOneBySlug, slug: %s, payload: %#v, Reason: %v", slug, d, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return ar, nil
}

func (s *ArticleService) GetArticles(ctx context.Context, args *conduit.ArticleArgs) ([]*model.Article, *ServiceError) {
	articles, err := s.articleRepo.Find(ctx, args)
	if err != nil {
		s.logger.Printf("Cannot Find::ArticleRepo args: %#v, Reason: %v", args, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	for _, a := range articles {
		tags, err := s.tagsRepo.GetArticleTagsByID(ctx, a.ID)
		if err != nil {
			s.logger.Printf("Cannot GetArticleTagsByID::ArticleTagsRepo for %s, Reason: %v", a.ID, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}

		a.TagList = tags

		if err := s.userRepo.FindOneByID(ctx, a.Author.ID, a.Author); err != nil {
			s.logger.Printf("Cannot FindOneByID::UserRepo for %s, Reason: %v", a.Author.ID, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return articles, nil
}

func (s *ArticleService) FavoriteArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *ServiceError) {
	a, err := s.GetOneBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if err := s.favoritesRepo.InsertOne(ctx, userID, a.ID); err != nil {
		s.logger.Printf("Cannot FavoriteArticle, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return a, nil
}

func (s *ArticleService) UnfavoriteArticleBySlug(ctx context.Context, userID, slug string) (*model.Article, *ServiceError) {
	a, err := s.GetOneBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if err := s.favoritesRepo.Delete(ctx, userID, a.ID); err != nil {
		s.logger.Printf("Cannot UnfavoriteArticle, Reason: %v", err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return a, nil
}

func (s *ArticleService) IsArticleFavorited(ctx context.Context, userID, articleID string) bool {
	return s.favoritesRepo.GetOneByIDs(ctx, userID, articleID) == nil
}
