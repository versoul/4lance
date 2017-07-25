package routes

import (
	"fmt"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

var (
	projectsPerPage = 50
)

func dashboardRoutes(r chi.Router) {
	r.Get("/dashboard/:page", dashboardPage)
}

func getUserData(userEmail string) map[string]interface{} {
	var result map[string]interface{}
	mgoSession, err := mgo.Dial(conf.DbHost)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	db := mgoSession.DB(conf.DbName).C("users")

	err = db.Find(bson.M{"email": userEmail}).One(&result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

func dashboardPage(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now().UnixNano() / 1000000
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	sessUser := sess.Get("user")
	fmt.Println("SESS")
	fmt.Println(sessUser)
	var userData = map[string]interface{}{}
	var userFilter = map[string]interface{}{}
	var userFilterCategories = []interface{}{}
	var userFilterKeywords = []interface{}{}
	if sessUser != nil {
		userData = getUserData(sessUser.(string))
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
	var timeProcess1 int64
	var timeProcess2 int64

	db := session.DB("4lance").C("projects")
	// Query All
	var results []map[string]interface{}
	if sessUser != nil {
		var query = bson.M{}
		bsonCategories := bson.M{"projectCategories": bson.M{"$in": userFilterCategories}}
		if len(userFilterKeywords) > 0 {
			var bsonKeywords = []bson.M{}
			for _, keyword := range userFilterKeywords {
				bsonKeywords = append(bsonKeywords, bson.M{"projectTitle": bson.M{"$regex": bson.RegEx{`\s` + keyword.(string) + `\s`, "gmi"}}})
				bsonKeywords = append(bsonKeywords, bson.M{"projectDescription": bson.M{"$regex": bson.RegEx{keyword.(string), "gmi"}}})
			}
			bsonOrKeywords := bson.M{"$or": bsonKeywords}
			query = bson.M{"$and": []bson.M{bsonCategories, bsonOrKeywords}}
		} else {
			query = bsonCategories
		}

		err = db.Find(query).Sort("-projectDate").Skip((page - 1) * projectsPerPage).Limit(projectsPerPage).All(&results)
		timeProcess1 = time.Now().UnixNano()/1000000 - timeStart
		m := map[string]interface{}{}
		err = db.Find(query).Sort("-projectDate").Skip((page - 1) * projectsPerPage).Limit(projectsPerPage).Explain(m)
		if err == nil {
			fmt.Printf("Explain: milis=%v docsExamined=%v \n", m["executionStats"].(map[string]interface{})["executionTimeMillis"].(int),
				m["executionStats"].(map[string]interface{})["totalDocsExamined"].(int))
		}
		timeProcess2 = time.Now().UnixNano()/1000000 - timeStart - timeProcess1
	} else {
		err = db.Find(nil).Sort("-projectDate").Skip((page - 1) * projectsPerPage).Limit(projectsPerPage).All(&results)
		m := map[string]interface{}{}
		err = db.Find(nil).Sort("-projectDate").Skip((page - 1) * projectsPerPage).Limit(projectsPerPage).Explain(m)
		if err == nil {
			fmt.Printf("Explain: milis=%v docsExamined=%v \n", m["executionStats"].(map[string]interface{})["executionTimeMillis"].(int),
				m["executionStats"].(map[string]interface{})["totalDocsExamined"].(int))
		}
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

	timeProcess3 := time.Now().UnixNano()/1000000 - timeStart - timeProcess2 - timeProcess1

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
	timeProcess4 := time.Now().UnixNano()/1000000 - timeStart - timeProcess3 - timeProcess2 - timeProcess1

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

	timeProcess1Str := strconv.FormatInt(timeProcess1, 10)
	fmt.Println("Before main query time = " + timeProcess1Str + "ms")
	timeProcess2Str := strconv.FormatInt(timeProcess2, 10)
	fmt.Println("After main query time = " + timeProcess2Str + "ms")
	timeProcess3Str := strconv.FormatInt(timeProcess3, 10)
	fmt.Println("Before loop time = " + timeProcess3Str + "ms")
	timeProcess4Str := strconv.FormatInt(timeProcess4, 10)
	fmt.Println("After loop time = " + timeProcess4Str + "ms")

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
