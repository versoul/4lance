package routes

import (
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"net/http"
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

		var userData = map[string]interface{}{}
		data, authOk := a.CheckAuthReq(r)
		if !authOk {
			http.Error(w, http.StatusText(403), 403)
			return
		} else {
			userData = data
			if userData["admin"].(bool) != true {
				http.Error(w, http.StatusText(403), 403)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func adminUsersPage(w http.ResponseWriter, r *http.Request) {
	var userData = map[string]interface{}{}
	data, authOk := a.CheckAuthReq(r)
	if authOk {
		userData = data
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//Get data

	db := session.DB("4lance").C("users")
	result := []map[string]interface{}{}
	err = db.Find(nil).All(&result)
	if err != nil {
		panic(err)
	}
	renderTemplate(w, "adminUsers", map[string]interface{}{
		"users": result,
		"user":  userData,
	})
}
func adminMainPage(w http.ResponseWriter, r *http.Request) {
	var userData = map[string]interface{}{}
	data, authOk := a.CheckAuthReq(r)
	if authOk {
		userData = data
	}

	renderTemplate(w, "adminMain", map[string]interface{}{
		"data": "<a href='/admin/users/'>Пользователи</a>",
		"user": userData,
	})
}
