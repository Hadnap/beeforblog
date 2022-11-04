package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Hadnap/beeforblog/api"
	"github.com/Hadnap/beeforblog/repository"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/russross/blackfriday"
)

var db *sql.DB

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

func viewHandler(w http.ResponseWriter, r *http.Request, slug string) {
	articleRepo := repository.NewStorage(db)
	articleService := api.NewArticleService(articleRepo)
	article, _ := articleService.Get(slug)
	renderTemplate(w, "view", article)
}

func editHandler(w http.ResponseWriter, r *http.Request, slug string) {
	articleRepo := repository.NewStorage(db)
	articleService := api.NewArticleService(articleRepo)
	article, _ := articleService.Get(slug)
	renderTemplate(w, "edit", article)
}

func saveHandler(w http.ResponseWriter, r *http.Request, slug string) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	article := api.ArticleRequest{Title: title, Content: content}
	articleRepo := repository.NewStorage(db)
	articleService := api.NewArticleService(articleRepo)
	slug, err := articleService.New(article)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+slug, http.StatusFound)
}

var viewTemplate = template.Must(template.New("view").Funcs(template.FuncMap{"markdown": markDowner}).ParseFiles("templates/index.html", "templates/view.html"))
var editTemplate = template.Must(template.New("edit").ParseFiles("templates/index.html", "templates/edit.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, article api.ArticleResponse) {
	var err error
	switch tmpl {
	case "view":
		err = viewTemplate.ExecuteTemplate(w, "base", article)
	case "edit":
		err = editTemplate.ExecuteTemplate(w, "base", article)
	default:
		err = viewTemplate.ExecuteTemplate(w, tmpl, article)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9\\-_]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func execMigrations(db *sql.DB) error {
	driver, _ := sqlite3.WithInstance(db, &sqlite3.Config{})

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver,
	)

	if err != nil {
		fmt.Printf("migration err: %v", err)
	}

	err = m.Up()

	switch err {
	case errors.New("no change"):
		return nil
	}

	return nil
}

func main() {
	connString := "file:main.db?cache=shared&mode=rwc"
	var err error
	db, err = sql.Open("sqlite3", connString)
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	err = execMigrations(db)

	if err != nil {
		fmt.Printf("err: %v", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	fileServer := http.FileServer(http.Dir("./static/"))

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
