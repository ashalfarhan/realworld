package controller

import (
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/test"
	"github.com/stretchr/testify/assert"
)

var articleController = &ArticleController{}

func TestCreateArticle(t *testing.T) {
	t.Parallel()

	t.Run("should return validation error in body, title, description", func(t *testing.T) {
		payload := response.M{
			"article": response.M{},
		}
		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/articles", payload, &body, articleController.CreateArticle)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		for k, v := range body {
			assert.Greaterf(t, len(v), 0, "validation error in %s", k)
		}
	})

	t.Run("tagList shoud be unique", func(t *testing.T) {
		payload := response.M{
			"article": response.M{
				"title":       "Article 1",
				"body":        "Article 1",
				"description": "Article 1",
				"tagList":     []string{"react", "react"},
			},
		}
		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/articles", payload, &body, articleController.CreateArticle)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Greater(t, len(body["errors"]["taglist"]), 0, "validation error in tagList")
		assert.Contains(t, body["errors"]["taglist"][0], "unique", "should contains unique message")
	})
}

func TestGetArticleList(t *testing.T) {
	t.Parallel()

	t.Run("Should response error if \"limit\" below 1", func(t *testing.T) {
		payload := response.M{}
		var body DtoError
		res := test.MakeRequest(http.MethodGet, "/api/articles?limit=-29", payload, &body, articleController.GetFiltered)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Greater(t, len(body["errors"]["limit"]), 0, "validation error in limit")
		assert.Contains(t, body["errors"]["limit"][0], "min 1")
	})

	t.Run("Should response error if \"limit\" above 25", func(t *testing.T) {
		payload := response.M{}
		var body DtoError
		res := test.MakeRequest(http.MethodGet, "/api/articles?limit=29", payload, &body, articleController.GetFiltered)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Greater(t, len(body["errors"]["limit"]), 0, "validation error in limit")
		assert.Contains(t, body["errors"]["limit"][0], "max 25")
	})

	t.Run("Should response error if \"offset\" below 0", func(t *testing.T) {
		payload := response.M{}
		var body DtoError
		res := test.MakeRequest(http.MethodGet, "/api/articles?offset=-29", payload, &body, articleController.GetFiltered)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Greater(t, len(body["errors"]["offset"]), 0, "validation error in offset")
		assert.Contains(t, body["errors"]["offset"][0], "min 0")
	})
}
