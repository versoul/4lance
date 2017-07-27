package routes

import (
	"encoding/json"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"versoul/4lance/auth"
	"versoul/4lance/config"
)

var (
	conf = config.GetInstance()
	a    = auth.GetInstance()
)

func init() {
	a.Configure("mongodb", map[string]interface{}{
		"dbHost":   conf.DbHost,
		"dbName":   conf.DbName,
		"mailHost": conf.MailHost,
		"mailPort": conf.MailPort,
		"mailUser": conf.MailUser,
		"mailPass": conf.MailPass,
	})
}

func authRoutes(r chi.Router) {
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
		"Type": "info",
		"Message": "На ваш email выслано письмо для активации аккаунта. Перейдите по ссылке в нем." +
			"<br/>Письмо может попасть в папку \"спам\" по ошибке, проверьте там.",
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

		mgoSession, err := mgo.Dial(conf.DbHost)
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		dbC := mgoSession.DB(conf.DbName).C("categories")
		// Query All
		var categories []map[string]interface{}
		err = dbC.Find(nil).Select(bson.M{"_id": 0}).All(&categories)
		if err != nil {
			panic(err)
		}

		var allCategories = []string{}
		for _, row := range categories[0] {
			row := row.([]interface{})
			for _, c := range row {
				cr := c.(map[string]interface{})
				childs := cr["childs"].([]interface{})
				for _, child := range childs {
					child1 := child.(map[string]interface{})
					tid := child1["tid"].(string)
					allCategories = append(allCategories, tid)
				}
			}
		}

		db := mgoSession.DB(conf.DbName).C("users")

		query := bson.M{"_id": id}
		change := bson.M{"$set": bson.M{
			"filter": bson.M{
				"categories": allCategories,
				"keywords":   []string{},
			},
		},
		}
		err = db.Update(query, change)
		w.Write([]byte("{\"status\":\"ok\"}"))
	}
}
func loginAction(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	_, err = a.LoginUser(w, data["email"].(string), data["password"].(string))
	if err != nil {
		w.Write([]byte("{\"status\":\"err\", \"error\":\"" + err.Error() + "\"}"))
	} else {
		w.Write([]byte("{\"status\":\"ok\"}"))
	}
}
func logoutAction(w http.ResponseWriter, r *http.Request) {
	a.LogoutUser(w)
	http.Redirect(w, r, "/dashboard/", 302)
}
