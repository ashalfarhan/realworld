package controller

import (
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/test"
)

var articleController = &ArticleController{}

func TestCreateArticle(t *testing.T) {
	t.Parallel()

	t.Run("Payload should have body, title, description", func(t *testing.T) {
		payload := response.M{
			"article": response.M{},
		}
		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/articles", payload, articleController.CreateArticle, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if len(body["errors"]["body"]) == 0 && len(body["errors"]["title"]) == 0 && len(body["errors"]["description"]) == 0 {
			t.Fatalf("expected error validation in body, title, and description. but got %v", body["errors"])
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
		res := test.MakeRequest(http.MethodPost, "/api/articles", payload, articleController.CreateArticle, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if len(body["errors"]["taglist"]) == 0 && body["errors"]["taglist"][0] != "unique " {
			t.Fatalf("expected validation errors in taglist, but got %v", body["errors"]["taglist"])
		}
	})
}

func TestGetArticleList(t *testing.T) {
	t.Parallel()

	t.Run("Should response error if \"limit\" below 1", func(t *testing.T) {
		payload := response.M{}
		var body DtoError
		res := test.MakeRequest(http.MethodGet, "/api/articles?limit=-29", payload, articleController.GetFiltered, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}
		if len(body["errors"]["limit"]) == 0 && body["errors"]["limit"][0] != "min 1" {
			t.Fatalf("expected validation errors in \"limit\", but got %v", body["errors"]["limit"])
		}
	})

	t.Run("Should response error if \"limit\" above 25", func(t *testing.T) {
		payload := response.M{}
		var body DtoError
		res := test.MakeRequest(http.MethodGet, "/api/articles?limit=29", payload, articleController.GetFiltered, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}
		if len(body["errors"]["limit"]) == 0 && body["errors"]["limit"][0] != "max 25" {
			t.Fatalf("expected validation errors in \"limit\", but got %v", body["errors"]["limit"])
		}
	})

	t.Run("Should response error if \"offset\" below 0", func(t *testing.T) {
		payload := response.M{}
		var body DtoError
		res := test.MakeRequest(http.MethodGet, "/api/articles?offset=-29", payload, articleController.GetFiltered, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}
		if len(body["errors"]["offset"]) == 0 && body["errors"]["offset"][0] != "min 0" {
			t.Fatalf("expected validation errors in \"offset\", but got %v", body["errors"]["offset"])
		}
	})
}
