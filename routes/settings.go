package routes

import (
	"encoding/json"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// Routes creates a REST router for the todos resource
/*func (rs dashboardResource) SettingsRoutes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Post("/saveFilter", rs.saveFilter) // GET /todos - read a list of todos
	return r
}*/
type Data struct {
	Categories []string `json'categories'`
	Keywords   []string `json'keywords'`
}

func settingsRoutes(r chi.Router) {
	r.Post("/filterSave", filterSave)
	r.Post("/filterReset", filterReset)
}

func filterSave(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	data := Data{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	sessUser := sess.Get("user")

	if sessUser == nil {
		w.Write([]byte("{\"status\":\"err\", \"error\":\"Нет! Нет! Только зарегистрированные!\"}"))
	} else {
		mgoSession, err := mgo.Dial(conf.DbHost)
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		db := mgoSession.DB(conf.DbName).C("users")

		query := bson.M{"email": sessUser}
		change := bson.M{"$set": bson.M{
			"filter": bson.M{
				"categories": data.Categories,
				"keywords":   data.Keywords,
			},
		},
		}
		err = db.Update(query, change)
		if err != nil {
			w.Write([]byte("{\"status\":\"err\", \"error\":\"" + err.Error() + "\"}"))
		} else {
			w.Write([]byte("{\"status\":\"ok\"}"))
		}
	}
}
func filterReset(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	sessUser := sess.Get("user")

	if sessUser == nil {
		w.Write([]byte("{\"status\":\"err\", \"error\":\"Нет! Нет! Только зарегистрированные!\"}"))
	} else {
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

		query := bson.M{"email": sessUser}
		change := bson.M{"$set": bson.M{
			"filter": bson.M{
				"categories": allCategories,
				"keywords":   []string{},
			},
		},
		}
		err = db.Update(query, change)
		if err != nil {
			w.Write([]byte("{\"status\":\"err\", \"error\":\"" + err.Error() + "\"}"))
		} else {
			w.Write([]byte("{\"status\":\"ok\"}"))
		}
	}
}
