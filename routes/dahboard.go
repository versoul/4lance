package routes

import (
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var (
	projectsPerPage = 100
)

type dashboardResource struct{}

// Routes creates a REST router for the todos resource
func (rs dashboardResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/:page", rs.List) // GET /todos - read a list of todos
	return r
}

func (rs dashboardResource) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil || page < 1 {
		page = 1
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	db := session.DB("4lance").C("projects")
	// Query All
	var results []map[string]interface{}
	err = db.Find(nil).Sort("-projectDate").Skip((page - 1) * projectsPerPage).Limit(projectsPerPage).All(&results)
	if err != nil {
		panic(err)
	}

	for _, v := range results {
		if str, ok := v["projectTitle"].(string); ok {
			v["projectTitle"] = template.HTML(str)
		}
		if str, ok := v["projectPrice"].(string); ok {
			v["projectPrice"] = template.HTML(str)
		}
		if str, ok := v["site"].(string); ok {
			v["projectIcon"] = hlprs.SiteToIcon(str)
		}
		href, ok1 := v["projectHref"].(string)
		site, ok2 := v["site"].(string)
		if ok1 && ok2 {
			v["projectHref"] = hlprs.ToFullLink(site, href)
		}
		if date, ok := v["projectDate"].(time.Time); ok {
			v["projectDate"] = hlprs.FormatTime(date, "02.01 | 15:04")
		}
	}

	cnt, err := db.Find(nil).Count()
	if err != nil {
		panic(err)
	}

	var pages []map[string]interface{}

	for i := page; i < page+10; i++ {
		p := map[string]interface{}{
			"num": i,
		}
		if page == i {
			p["active"] = true
		}
		if cnt-i*projectsPerPage > 0 {
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

	renderTemplate(w, "main", map[string]interface{}{
		"Error":      "Main Website",
		"projects":   results,
		"count":      cnt,
		"pagination": pagination,
	})
}
