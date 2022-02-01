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
	articleRepo     *repository.ArticleRepository
	userRepo        *repository.UserRepository
	articleTagsRepo *repository.ArticleTagsRepository
	logger          *log.Logger
}

func NewArticleService(repo *repository.Repository) *ArticleService {
	return &ArticleService{
		articleRepo:     repo.ArticleRepo,
		userRepo:        repo.UserRepo,
		articleTagsRepo: repo.ArticleTagsRepo,
		logger:          conduit.NewLogger("article-service"),
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
			if err := s.articleTagsRepo.InsertOne(ctx, a.ID, tag); err != nil {
				s.logger.Printf("Cannot InsertOne ArticleTags Repo for %s, Reason: %v", tag, err)
				return nil, CreateServiceError(http.StatusInternalServerError, nil)
			}

			a.TagList = append(a.TagList, tag)
		}
	}

	if err := s.userRepo.FindOneById(ctx, a.Author.ID, a.Author); err != nil {
		s.logger.Printf("Cannot FindOneById User Repo for %s, Reason: %v", a.Author.ID, err)
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

	tags, err := s.articleTagsRepo.GetArticleTagsById(ctx, a.ID)
	if err != nil {
		s.logger.Printf("Cannot GetArticleTagsById ArticleTags Repo for %s, Reason: %v", a.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.TagList = tags

	if err := s.userRepo.FindOneById(ctx, a.Author.ID, a.Author); err != nil {
		s.logger.Printf("Cannot FindOneById User Repo for %s, Reason: %v", a.Author.ID, err)
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
	tags, err := s.articleTagsRepo.GetAllTags(ctx)
	if err != nil {
		return nil, CreateServiceError(http.StatusBadRequest, err)
	}

	return tags, nil
}

func (s *ArticleService) CreateSlug(title string) string {
	id, err := gonanoid.New(defaultSlugId)
	if err != nil {
		s.logger.Printf("Cannot create nanoid for %s, %v", title, err)
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

	if err := s.articleRepo.Update(ctx, slug, args, ar); err != nil {
		s.logger.Printf("Cannot update article, slug:%s, payload:%#v, Reason:%v", slug, d, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return ar, nil
}

func (s *ArticleService) GetArticles(ctx context.Context, args *conduit.ArticleArgs) ([]model.Article, *ServiceError) {
	articles, err := s.articleRepo.Find(ctx, args)
	if err != nil {
		s.logger.Printf("Cannot Find args:%v, Reason:%v", args, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	for _, a := range articles {
		tags, err := s.articleTagsRepo.GetArticleTagsById(ctx, a.ID)
		if err != nil {
			s.logger.Printf("Cannot GetArticleTagsById ArticleTags Repo for %s, Reason: %v", a.ID, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}

		a.TagList = tags

		if err := s.userRepo.FindOneById(ctx, a.Author.ID, a.Author); err != nil {
			s.logger.Printf("Cannot FindOneById User Repo for %s, Reason: %v", a.Author.ID, err)
			return nil, CreateServiceError(http.StatusInternalServerError, nil)
		}
	}

	return articles, nil
}
