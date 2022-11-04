package repository

import (
	"database/sql"
	"github.com/gosimple/slug"
	_ "github.com/mattn/go-sqlite3"
	"log"

	"github.com/Hadnap/beeforblog/api"
)

type Storage interface {
	GetArticleBySlug(slug string) (api.ArticleResponse, error)
	CreateArticle(request api.ArticleRequest) (string, error)
	UpdateArticle(slugId string, request api.ArticleRequest) (string, error)
}

type storage struct {
	db *sql.DB
}

func (s *storage) GetArticleBySlug(slug string) (api.ArticleResponse, error) {
	getArticleStatement := "SELECT title, content, created_at FROM articles WHERE slug = $1"
	var (
		title     string
		content   string
		createdAt string
	)
	err := s.db.QueryRow(getArticleStatement, slug).Scan(&title, &content, &createdAt)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return api.ArticleResponse{}, err
	}

	return api.ArticleResponse{
		Title:     title,
		Content:   content,
		Slug:      slug,
		CreatedAt: createdAt,
	}, nil
}

func (s *storage) CreateArticle(request api.ArticleRequest) (string, error) {
	sluggedTitle := slug.Make(request.Title)
	newArticleStatement := `
		INSERT INTO articles (title, slug, content) VALUES ($1, $2, $3)
		`
	err := s.db.QueryRow(newArticleStatement, request.Title, sluggedTitle, request.Content).Err()
	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return "", err
	}

	return sluggedTitle, nil
}

func (s *storage) UpdateArticle(slugId string, request api.ArticleRequest) (string, error) {
	sluggedTitle := slug.Make(request.Title)
	updateArticleStatement := `
		UPDATE articles SET (title, slug, content) VALUES ($1, $2, $3) WHERE slug = $4
		`
	err := s.db.QueryRow(updateArticleStatement, request.Title, sluggedTitle, request.Content, slugId).Err()
	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return "", err
	}

	return sluggedTitle, nil
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}
