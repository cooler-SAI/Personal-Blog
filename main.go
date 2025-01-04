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
	// Гостевые маршруты
	http.HandleFunc("/", homePage)
	http.HandleFunc("/article/", articlePage)

	// Админские маршруты с защитой
	http.HandleFunc("/admin", adminMiddleware(adminDashboard))
	http.HandleFunc("/admin/new", adminMiddleware(newArticlePage))
	http.HandleFunc("/admin/edit/", adminMiddleware(editArticlePage))
	http.HandleFunc("/admin/delete/", adminMiddleware(deleteArticle))

	// Статические файлы (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Запуск сервера
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// Обработчик главной страницы
	templates.ExecuteTemplate(w, "home.html", nil)
}

func articlePage(w http.ResponseWriter, r *http.Request) {
	// Обработчик страницы статьи
	templates.ExecuteTemplate(w, "article.html", nil)
}

func adminDashboard(w http.ResponseWriter, r *http.Request) {
	// Пример данных для шаблона
	articles := []struct {
		ID    int
		Title string
	}{
		{ID: 1, Title: "My First Article"},
		{ID: 2, Title: "Second Article"},
	}

	// Рендеринг шаблона
	err := templates.ExecuteTemplate(w, "dashboard.html", articles)
	if err != nil {
		return
	}
}

func newArticlePage(w http.ResponseWriter, r *http.Request) {
	// Страница добавления новой статьи
	templates.ExecuteTemplate(w, "new.html", nil)
}

func editArticlePage(w http.ResponseWriter, r *http.Request) {
	// Страница редактирования статьи
	templates.ExecuteTemplate(w, "edit.html", nil)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// Удаление статьи (заглушка)
	_, err := w.Write([]byte("Article deleted!"))
	if err != nil {
		return
	}
}
