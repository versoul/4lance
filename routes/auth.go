package routes

import (
	"fmt"
	"github.com/pressly/chi"
	"net/http"
	"versoul/4lance/auth"
	"versoul/4lance/config"
)

var (
	conf = config.GetInstance()
	a    = auth.GetInstance("mongodb", map[string]interface{}{
		"dbHost":   "localhost",
		"dbName":   "4lance",
		"mailHost": conf.MailHost,
		"mailPort": conf.MailPort,
		"mailUser": conf.MailUser,
		"mailPass": conf.MailPass,
	})
)

func registerPage(w http.ResponseWriter, r *http.Request) {
	id, err := a.RegisterUser("mailto.versoul@gmail.com", "111")
	if err != nil {
		fmt.Println("Register fail - " + err.Error())
	} else {
		fmt.Println("Register success")
		fmt.Println(id)
		a.SendConfirmationEmail(id)
	}
	renderTemplate(w, "register", map[string]interface{}{})
}
func confirmEmailPage(w http.ResponseWriter, r *http.Request) {
	err := a.ConfirmEmail(chi.URLParam(r, "confirmationHash"))
	if err != nil {
		renderTemplate(w, "message", map[string]interface{}{
			"Type":    "danger",
			"Message": err.Error(),
		})
	} else {
		renderTemplate(w, "message", map[string]interface{}{
			"Type":    "success",
			"Message": "Email подтвержден и аккаунт успешно активирован. Можете авторизироваться и пользоваться сервисом.",
		})
	}

}
func loginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", map[string]interface{}{})
}
