package main

import (
	"html/template"
	"net/http"
)

var templates *template.Template

func init() {
	// Загружаем шаблоны
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

// Middleware для авторизации администратора
func adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка базовой авторизации
		username, password, ok := r.BasicAuth()
		if !ok || username != "admin" || password != "password" {
			// Если авторизация не удалась
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Если авторизация успешна, продолжаем выполнение следующего обработчика
		next(w, r)
	}
}

func main() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/article/", articlePage)

	http.HandleFunc("/admin", adminMiddleware(adminDashboard))
	http.HandleFunc("/admin/new", adminMiddleware(newArticlePage))
	http.HandleFunc("/admin/edit/", adminMiddleware(editArticlePage))
	http.HandleFunc("/admin/delete/", adminMiddleware(deleteArticle))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func homePage(w http.ResponseWriter, _ *http.Request) {

	err := templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		return
	}
}

func articlePage(w http.ResponseWriter, _ *http.Request) {

	err := templates.ExecuteTemplate(w, "article.html", nil)
	if err != nil {
		return
	}
}

func adminDashboard(w http.ResponseWriter, _ *http.Request) {

	articles := []struct {
		ID    int
		Title string
	}{
		{ID: 1, Title: "My First Article"},
		{ID: 2, Title: "Second Article"},
	}

	err := templates.ExecuteTemplate(w, "dashboard.html", articles)
	if err != nil {
		return
	}
}

func newArticlePage(w http.ResponseWriter, _ *http.Request) {

	err := templates.ExecuteTemplate(w, "new.html", nil)
	if err != nil {
		return
	}
}

func editArticlePage(w http.ResponseWriter, _ *http.Request) {

	err := templates.ExecuteTemplate(w, "edit.html", nil)
	if err != nil {
		return
	}
}

func deleteArticle(w http.ResponseWriter, _ *http.Request) {

	_, err := w.Write([]byte("Article deleted!"))
	if err != nil {
		return
	}
}
