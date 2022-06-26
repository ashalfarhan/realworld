package service_test

import (
	"testing"

	"github.com/ashalfarhan/realworld/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateArticle(t *testing.T) {
	authorID := ""
	u := &model.User{
		ID: authorID,
	}
	userRepoMock.On("FindOneByID", mockCtx, authorID).Return(u, nil).Once()
	as := assert.New(t)

	d := &model.CreateArticleFields{
		Title:   "My first article",
		TagList: []string{"typescript", "react", "javascript", "golang"},
	}

	articleTagsRepoMock.On("InsertBulk", mockCtx, mock.Anything).Return(nil).Once()
	articleRepoMock.On("InsertOne", mockCtx, d, authorID).Return(&model.Article{}, nil).Once()
	a, err := articleService.CreateArticle(tctx, d, authorID)
	articleRepoMock.AssertExpectations(t)
	userRepoMock.AssertExpectations(t)
	articleTagsRepoMock.AssertExpectations(t)

	as.Nil(err)
	as.NotNil(a)
	as.Equal(d.TagList, a.TagList, "Tag list should be set")
	as.NotEqual(d.Slug, d.Title, "Slug should be different from title")
	as.Greater(d.Slug, d.Title, "Slug length must be greater than title, and added id")
}
