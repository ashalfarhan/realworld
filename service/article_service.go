package service

import (
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

func (s *ArticleService) Create(d *dto.CreateArticleDto, authorID string) (*model.Article, *ServiceError) {
	a := &model.Article{
		Title:       d.Title,
		Description: d.Description,
		Body:        d.Body,
		Author: &model.User{
			ID: authorID,
		},
	}

	a.Slug = s.CreateSlug(a.Title)

	if err := s.articleRepo.InsertOne(a); err != nil {
		s.logger.Printf("Cannot InsertOne to ArticleRepo for %#v, Reason: %v", a, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	if len(d.TagList) > 0 {
		for _, tag := range d.TagList {
			if err := s.articleTagsRepo.InsertOne(a.ID, tag); err != nil {
				s.logger.Printf("Cannot InsertOne ArticleTags Repo for %s, Reason: %v", tag, err)
				return nil, CreateServiceError(http.StatusInternalServerError, nil)
			}

			a.TagList = append(a.TagList, tag)
		}
	}

	if err := s.userRepo.FindOneById(a.Author.ID, a.Author); err != nil {
		s.logger.Printf("Cannot FindOneById User Repo for %s, Reason: %v", a.Author.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return a, nil
}

func (s *ArticleService) GetOneBySlug(slug string) (*model.Article, *ServiceError) {
	a := &model.Article{
		Slug:   slug,
		Author: &model.User{},
	}

	if err := s.articleRepo.FindOneBySlug(a); err != nil {
		if err == sql.ErrNoRows {
			return nil, CreateServiceError(http.StatusNotFound, errors.New("no article found"))
		}
		s.logger.Printf("Cannot FindOneBySlug Article Repo for %#v, Reason: %v", a, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	tags, err := s.articleTagsRepo.GetArticleTagsById(a.ID)
	if err != nil {
		s.logger.Printf("Cannot GetArticleTagsById ArticleTags Repo for %s, Reason: %v", a.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	a.TagList = tags

	if err := s.userRepo.FindOneById(a.Author.ID, a.Author); err != nil {
		s.logger.Printf("Cannot FindOneById User Repo for %s, Reason: %v", a.Author.ID, err)
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	return a, nil
}

func (s *ArticleService) DeleteArticle(slug string, userID string) *ServiceError {
	a, err := s.GetOneBySlug(slug)
	if err != nil {
		return err
	}

	if a.Author.ID != userID {
		return CreateServiceError(http.StatusForbidden, errors.New("you cannot delete this article"))
	}

	if err := s.articleRepo.DeleteBySlug(slug); err != nil {
		s.logger.Printf("Cannot DeleteArticleBySlug Article Repo for %s, Reason: %v", slug, err)
		return CreateServiceError(http.StatusInternalServerError, nil)
	}

	return nil
}

func (s *ArticleService) GetAllTags() ([]string, *ServiceError) {
	tags, err := s.articleTagsRepo.GetAllTags()
	if err != nil {
		return nil, CreateServiceError(http.StatusBadRequest, err)
	}

	return tags, nil
}

func (s *ArticleService) CreateSlug(title string) string {
	id, err := gonanoid.New(defaultSlugId)
	if err != nil {
		log.Printf("[Slug Generation Error] Cannot create slug for %s, %s\n", title, err.Error())
	}

	slug := slug.Make(title)

	return fmt.Sprintf("%s-%s", slug, id)
}
