package routes

import (
	"github.com/pressly/chi"
	"net/http"
)

func projectRoutes(r chi.Router) {
	r.Get("/project/{id}", projectPage)
}

func projectPage(w http.ResponseWriter, r *http.Request) {
	var userData = map[string]interface{}{}
	//var userFilter = map[string]interface{}{}
	//var userFilterCategories = []interface{}{}
	//var userFilterKeywords = []interface{}{}
	data, authOk := a.CheckAuthReq(r)
	if authOk {
		userData = data
		//userFilter = userData["filter"].(map[string]interface{})
		//userFilterCategories = userFilter["categories"].([]interface{})
		//userFilterKeywords = userFilter["keywords"].([]interface{})
	}

	renderTemplate(w, "project", map[string]interface{}{
		"user": userData,
	})
}
