package routes

import (
	"fmt"
	"github.com/pressly/chi"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"versoul/4lance/templateHelpers"
)

var (
	hlprs = templateHelpers.GetHelpers()
)

func minus(a, b int64) string {
	fmt.Println(a)
	fmt.Println(b)
	return "12"
	//return strconv.FormatInt(a-b, 10)
}

func renderTemplate(w http.ResponseWriter, name string, pageData map[string]interface{}) {
	w.Header().Set("Content-type", "text/html")

	var funcMap = template.FuncMap{
		"minus": minus,
		"inc": func(i int) int {
			return i + 1
		},
	}

	var tpl = template.Must(template.New(name).Funcs(funcMap).ParseFiles(
		"./templates/base.html",
		"./templates/header.html",
		"./templates/"+name+".html"))
	tpl.ExecuteTemplate(w, "base", pageData)
}
func InitRoutes() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard/", 302)
	})
	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard/", 302)
	})
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	r.FileServer("/static", http.Dir(filesDir))

	r.Mount("/dashboard/", dashboardResource{}.Routes())

	http.ListenAndServe(":8080", r)
}
