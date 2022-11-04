package api

type ArticleRequest struct {
	Title   string
	Content string
}

type ArticleResponse struct {
	Title     string
	Slug      string
	Content   string
	CreatedAt string
}
