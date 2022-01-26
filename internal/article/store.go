package article

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/SergeyParamoshkin/rest/internal/model"
)

// Article fixture data
// nolint
var articles = []*model.Article{
	{ID: "1", UserID: 100, Title: "Hi", Slug: "hi"},
	{ID: "2", UserID: 200, Title: "sup", Slug: "sup"},
	{ID: "3", UserID: 300, Title: "alo", Slug: "alo"},
	{ID: "4", UserID: 400, Title: "bonjour", Slug: "bonjour"},
	{ID: "5", UserID: 500, Title: "whats up", Slug: "whats-up"},
}

func dbNewArticle(article *model.Article) (string, error) {
	article.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	articles = append(articles, article)

	return article.ID, nil
}

func dbGetArticle(id string) (*model.Article, error) {
	for _, a := range articles {
		if a.ID == id {
			return a, nil
		}
	}

	return nil, errors.New("article not found.")
}
func dbGetArticleBySlug(slug string) (*model.Article, error) {
	for _, a := range articles {
		if a.Slug == slug {
			return a, nil
		}
	}

	return nil, errors.New("Article not found.")
}

func dbUpdateArticle(id string, article *model.Article) (*model.Article, error) {
	for i, a := range articles {
		if a.ID == id {
			articles[i] = article

			return article, nil
		}
	}

	return nil, errors.New("article not found.")
}

func dbRemoveArticle(id string) (*model.Article, error) {
	for i, a := range articles {
		if a.ID == id {
			articles = append((articles)[:i], (articles)[i+1:]...)

			return a, nil
		}
	}

	return nil, errors.New("article not found.")
}
