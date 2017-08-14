package routes

import (
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

var (
	projectsPerPage = 50
)

func dashboardRoutes(r chi.Router) {
	r.Route("/dashboard/", func(r chi.Router) {
		r.Get("/", dashboardPage)
		r.Get("/{page}", dashboardPage)
	})
}

func dashboardPage(w http.ResponseWriter, r *http.Request) {
	var userData = map[string]interface{}{}
	var userFilter = map[string]interface{}{}
	var userFilterCategories = []interface{}{}
	var userFilterKeywords = []interface{}{}
	data, authOk := a.CheckAuthReq(r)
	if authOk {
		userData = data
		userFilter = userData["filter"].(map[string]interface{})
		userFilterCategories = userFilter["categories"].([]interface{})
		userFilterKeywords = userFilter["keywords"].([]interface{})
	}

	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil || page < 1 {
		page = 1
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//Get data

	db := session.DB("4lance").C("projects")
	// Query All
	var results []map[string]interface{}
	if authOk {
		var query = bson.M{}
		bsonCategories := bson.M{"projectCategories": bson.M{"$in": userFilterCategories}}
		if len(userFilterKeywords) > 0 {
			var bsonKeywords = []bson.M{}
			for _, keyword := range userFilterKeywords {
				bsonKeywords = append(bsonKeywords, bson.M{"projectTitle": bson.M{"$regex": bson.RegEx{keyword.(string), "mi"}}})
				bsonKeywords = append(bsonKeywords, bson.M{"projectDescription": bson.M{"$regex": bson.RegEx{keyword.(string), "mi"}}})
			}
			bsonOrKeywords := bson.M{"$or": bsonKeywords}
			query = bson.M{"$and": []bson.M{bsonCategories, bsonOrKeywords}}
		} else {
			query = bsonCategories
		}

		err = db.Find(query).Sort("-projectDate").Skip((page - 1) * projectsPerPage).Limit(projectsPerPage).All(&results)
	} else {
		err = db.Find(nil).Sort("-projectDate").Skip((page - 1) * projectsPerPage).Limit(projectsPerPage).All(&results)
	}

	if err != nil {
		panic(err)
	}

	//timeProcess2 := time.Now().UnixNano()/1000000 - timeStart - timeProcess1

	dbC := session.DB("4lance").C("categories")
	// Query All
	var categories []map[string]interface{}
	err = dbC.Find(nil).All(&categories)
	if err != nil {
		panic(err)
	}

	//Сравниваем категории с выбранными в сессии, помечаем для интерфейса
	for _, row := range categories[0] {
		row, ok := row.([]interface{})
		if ok {
			for _, c := range row {
				cr := c.(map[string]interface{})
				childs := cr["childs"].([]interface{})
				for _, child := range childs {
					child1 := child.(map[string]interface{})
					tid := child1["tid"].(string)
					for _, v := range userFilterCategories {
						if v.(string) == tid {
							child1["activ"] = true
						}
					}
				}
			}
		}
	}

	//Get data count
	cnt, err := db.Find(nil).Count()
	if err != nil {
		panic(err)
	}

	//Create pagination
	var pages []map[string]interface{}
	for i := page; i < page+9; i++ {
		p := map[string]interface{}{
			"num": i,
		}
		if page == i {
			p["active"] = true
		}
		if cnt-(i-1)*projectsPerPage >= 0 {
			pages = append(pages, p)
		}
	}
	prev := page - 1
	if prev < 1 {
		prev = 1
	}
	next := pages[len(pages)-1]["num"].(int)
	if cnt-(next+1)*projectsPerPage > 0 {
		next++
	}
	pagination := map[string]interface{}{
		"pages": pages,
		"prev":  prev,
		"next":  next,
	}

	//Render template
	renderTemplate(w, "main", map[string]interface{}{
		"Error":              "Main Website",
		"projects":           results,
		"categories":         categories,
		"count":              cnt,
		"pagination":         pagination,
		"userFilterKeywords": userFilterKeywords,
		"user":               userData,
	})
}
