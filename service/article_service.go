package service

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
)

type ArticleService struct {
	articleRepo     *repository.ArticleRepository
	userRepo        *repository.UserRepository
	articleTagsRepo *repository.ArticleTagsRepository
}

func NewArticleService(repo *repository.Repository) *ArticleService {
	return &ArticleService{
		articleRepo:     repo.ArticleRepo,
		userRepo:        repo.UserRepo,
		articleTagsRepo: repo.ArticleTagsRepo,
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

	a.Slug = conduit.CreateSlug(a.Title)

	if err := s.articleRepo.InsertOne(a); err != nil {
		return nil, CreateServiceError(http.StatusBadRequest, err)
	}

	if len(d.TagList) > 0 {
		for _, tag := range d.TagList {
			if err := s.articleTagsRepo.InsertOne(a.ID, tag); err != nil {
				return nil, CreateServiceError(http.StatusBadRequest, err)
			}

			a.TagList = append(a.TagList, tag)
		}
	}

	if err := s.userRepo.FindOneById(a.Author.ID, a.Author); err != nil {
		return nil, CreateServiceError(http.StatusBadRequest, err)
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

		return nil, CreateServiceError(http.StatusBadRequest, err)
	}

	tags, err := s.articleTagsRepo.GetArticleTagsById(a.ID)
	if err != nil {
		return nil, CreateServiceError(http.StatusBadRequest, err)
	}

	a.TagList = tags

	if err := s.userRepo.FindOneById(a.Author.ID, a.Author); err != nil {
		return nil, CreateServiceError(http.StatusBadRequest, err)
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
		return CreateServiceError(http.StatusBadRequest, err)
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
