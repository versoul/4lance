package routes

import (
	"fmt"
	"github.com/pressly/chi"
	"html/template"
	"net/http"

	"github.com/astaxie/beego/session"
	"os"
	"path/filepath"
	"versoul/4lance/templateHelpers"
)

var (
	hlprs          = templateHelpers.GetHelpers()
	localIp        = "192.168.1.91"
	globalSessions *session.Manager
)
var EcternalIp string

func minus(a, b int64) string {
	fmt.Println(a)
	fmt.Println(b)
	return "12"
	//return strconv.FormatInt(a-b, 10)
}

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "sessionid",
		EnableSetCookie: true,
		Gclifetime:      86400,
		CookieLifeTime:  86400,
		Secure:          false,
		ProviderConfig:  "./tmp",
	}
	globalSessions, _ = session.NewManager("file", sessionConfig)
	go globalSessions.GC()
}

func renderTemplate(w http.ResponseWriter, name string, pageData map[string]interface{}) {
	w.Header().Set("Content-type", "text/html")

	var funcMap = template.FuncMap{
		"minus": minus,
		"inc": func(i int) int {
			return i + 1
		},
	}

	templates := []string{
		"./templates/base.html",
		"./templates/header.html",
	}

	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	//Подклчаем метрики только на реальном сервере
	if ip != localIp {
		templates = append(templates, "./templates/metrics.html")
	}
	//Подключаем запрошеный шаблон
	templates = append(templates, "./templates/"+name+".html")
	//Парсим все шаблоны
	var tpl = template.New(name).Funcs(funcMap)
	for i := 0; i < len(templates); i++ {
		tpl.ParseFiles(templates[i])
	}
	//Отображаем
	tpl.ExecuteTemplate(w, "base", pageData)
}
func InitRoutes() {
	r := chi.NewRouter()
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	r.FileServer("/static", http.Dir(filesDir))

	staticRoutes(r)

	authRoutes(r)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard/", http.StatusFound)
	})
	r.Mount("/dashboard/", dashboardResource{}.Routes())
	r.Post("/saveFilter", saveFilter)
	r.Post("/saveKeyWords", saveKeyWords)

	panic(http.ListenAndServe(":8080", r))
}
func staticRoutes(r chi.Router) {
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	})
	r.Get("/google173377f79f6a476a.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/google173377f79f6a476a.html")
	})
	r.Get("/yandex_87613e29f1d00477.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/yandex_87613e29f1d00477.html")
	})
}
