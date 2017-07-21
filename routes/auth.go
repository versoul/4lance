package routes

import (
	"encoding/json"
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

func init() {
	r.Get("/register/", registerPage)
	r.Get("/login/", loginPage)
	r.Get("/confirmMessage/", confirmMessagePage)
	r.Get("/confirmEmail/:confirmationHash/", confirmEmailPage)
	r.Get("/logout/", logoutAction)

	r.Post("/register/", registerAction)
	r.Post("/login/", loginAction)
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register", map[string]interface{}{})
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", map[string]interface{}{})
}
func confirmMessagePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "message", map[string]interface{}{
		"Type":    "info",
		"Message": "На ваш email выслано письмо для активации аккаунта. Перейдите по ссылке в нем.",
	})
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

func registerAction(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	id, err := a.RegisterUser(data["email"].(string), data["password"].(string))
	if err != nil {
		w.Write([]byte("{\"status\":\"err\", \"error\":\"" + err.Error() + "\"}"))
	} else {
		a.SendConfirmationEmail(id)
		w.Write([]byte("{\"status\":\"ok\"}"))
	}
}
func loginAction(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	userData, err := a.LoginUser(data["email"].(string), data["password"].(string))
	if err != nil {
		w.Write([]byte("{\"status\":\"err\", \"error\":\"" + err.Error() + "\"}"))
	} else {
		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)
		sess.Set("user", userData)
		w.Write([]byte("{\"status\":\"ok\"}"))
	}
}
func logoutAction(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sess.Flush()
	http.Redirect(w, r, "/dashboard/", 302)
}
