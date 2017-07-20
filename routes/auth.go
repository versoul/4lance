package routes

import (
	"fmt"
	"net/http"
	"versoul/4lance/auth"
)

var (
	a = auth.GetInstance("mongodb", map[string]interface{}{
		"host":   "localhost",
		"dbName": "4lance",
	})
)

func registerPage(w http.ResponseWriter, r *http.Request) {
	err := a.RegisterUser("E", "P")
	if err != nil {
		fmt.Println("Register fail" + err.Error())
	} else {
		fmt.Println("Register success")
	}
	renderTemplate(w, "register", map[string]interface{}{})
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", map[string]interface{}{})
}
