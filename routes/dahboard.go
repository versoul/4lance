package routes

import (
	"errors"
	"fmt"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"html/template"
	"net"
	"net/http"
	"strconv"
	//"time"
)

var (
	projectsPerPage = 50
)

type dashboardResource struct{}

func init() {

}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
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

// Routes creates a REST router for the todos resource
func (rs dashboardResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/:page", rs.List) // GET /todos - read a list of todos
	return r
}

func (rs dashboardResource) List(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	sessUser := sess.Get("user")
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
	db := session.DB("4lance").C("projects")
	// Query All
	var results []map[string]interface{}
	if sessUser != nil {
		var query = bson.M{}
		bsonCategories := bson.M{"projectCategories": bson.M{"$in": userFilterCategories}}
		if len(userFilterKeywords) > 0 {
			var bsonKeywords = []bson.M{}
			for _, keyword := range userFilterKeywords {
				bsonKeywords = append(bsonKeywords, bson.M{"projectTitle": bson.M{"$regex": bson.RegEx{`\s`+keyword.(string)+`\s`, "gmi"}}})
				bsonKeywords = append(bsonKeywords, bson.M{"projectDescription": bson.M{"$regex": bson.RegEx{keyword.(string), "gmi"}}})
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
