package api

import (
	"errors"
)

type ArticleService interface {
	New(article ArticleRequest) (string, error)
	Update(slug string, article ArticleRequest) (string, error)
	Get(slug string) (ArticleResponse, error)
}

type ArticleRepository interface {
	CreateArticle(ArticleRequest) (string, error)
	UpdateArticle(string, ArticleRequest) (string, error)
	GetArticleBySlug(slug string) (ArticleResponse, error)
}

type articleService struct {
	storage ArticleRepository
}

func (a *articleService) Get(slug string) (ArticleResponse, error) {
	return a.storage.GetArticleBySlug(slug)
}

func (a *articleService) New(article ArticleRequest) (string, error) {
	if article.Title == "" {
		return "", errors.New("article service: missing title")
	}
	slug, err := a.storage.CreateArticle(article)
	if err != nil {
		return "", err
	}

	return slug, nil
}

func NewArticleService(articleRepo ArticleRepository) ArticleService {
	return &articleService{
		storage: articleRepo,
	}

}
