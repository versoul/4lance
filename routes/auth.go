package routes

import (
	"net/http"
)

func registerPage(w http.ResponseWriter, r *http.Request) {
	//Render template
	renderTemplate(w, "register", map[string]interface{}{})
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	//Render template
	renderTemplate(w, "login", map[string]interface{}{})
}
