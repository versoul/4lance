package routes

import (
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"net/http"
	"strconv"
)

func adminRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(adminOnly)

	r.Get("/", adminMainPage)
	r.Get("/users/", adminUsersPage)

	return r
}

func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)

		sessUser := sess.Get("user")

		var userData = map[string]interface{}{}

		if sessUser == nil {
			http.Error(w, http.StatusText(403), 403)
			return
		} else {
			userData = getUserData(sessUser.(string))
			if userData["admin"].(bool) != true {
				http.Error(w, http.StatusText(403), 403)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func adminUsersPage(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sessUser := sess.Get("user")
	var userData = map[string]interface{}{}
	if sessUser != nil {
		userData = getUserData(sessUser.(string))
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//Get data

	db := session.DB("4lance").C("users")
	cnt, err := db.Find(nil).Count()
	if err != nil {
		panic(err)
	}
	renderTemplate(w, "message", map[string]interface{}{
		"Type":    "success",
		"Message": "Admin users page. Users count: " + strconv.Itoa(cnt),
		"user":    userData,
	})
}
func adminMainPage(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sessUser := sess.Get("user")
	var userData = map[string]interface{}{}
	if sessUser != nil {
		userData = getUserData(sessUser.(string))
	}

	renderTemplate(w, "message", map[string]interface{}{
		"Type":    "success",
		"Message": "Admin main page",
		"user":    userData,
	})
}
