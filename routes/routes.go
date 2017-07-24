package routes

import (
	"fmt"
	"github.com/pressly/chi"
	"html/template"
	"net/http"

	"errors"
	"github.com/astaxie/beego/session"
	"net"
	"os"
	"path/filepath"
	"versoul/4lance/templateHelpers"
)

var (
	localIp        = "192.168.1.91"
	globalSessions *session.Manager
)

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
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
	var tpl = template.New(name).Funcs(templateHelpers.Helpers)
	for i := 0; i < len(templates); i++ {
		tpl.ParseFiles(templates[i])
	}
	//Отображаем
	tpl.ExecuteTemplate(w, "base", pageData)
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

func InitRoutes() {
	r := chi.NewRouter()
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	r.FileServer("/static", http.Dir(filesDir))

	staticRoutes(r)
	authRoutes(r)
	settingsRoutes(r)
	dashboardRoutes(r)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard/", http.StatusFound)
	})

	panic(http.ListenAndServe(":8080", r))
}
