package service_test

import (
	"testing"

	"github.com/ashalfarhan/realworld/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateArticle(t *testing.T) {
	username := ""
	u := &model.User{
		Username: username,
	}
	userRepoMock.On("FindOneByUsername", mockCtx, username).Return(u, nil)
	as := assert.New(t)

	d := &model.CreateArticleFields{
		Title:   "My first article",
		TagList: []string{"typescript", "react", "javascript", "golang"},
	}

	articleTagsRepoMock.On("InsertBulk", mockCtx, mock.Anything).Return(nil)
	articleRepoMock.On("InsertOne", mockCtx, d, username).Return(&model.Article{}, nil)
	a, err := articleService.CreateArticle(tctx, d, username)
	articleRepoMock.AssertExpectations(t)
	userRepoMock.AssertExpectations(t)
	articleTagsRepoMock.AssertExpectations(t)

	as.Nil(err)
	as.NotNil(a)
	as.Equal(d.TagList, a.TagList, "Tag list should be set")
	as.NotEqual(d.Slug, d.Title, "Slug should be different from title")
	as.Greater(d.Slug, d.Title, "Slug length must be greater than title, and added id")
}
