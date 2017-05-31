package routes

import (
	"errors"
	"fmt"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	projectsPerPage = 100
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
	//Get data
	db := session.DB("4lance").C("projects")
	// Query All
	var results []map[string]interface{}
	err = db.Find(nil).Sort("-projectDate").Skip((page - 1) * projectsPerPage).Limit(projectsPerPage).All(&results)
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
	//Works with data
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
			v["projectDate"] = hlprs.FormatTime(date, "02.01 15:04")
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
		"Error":      "Main Website",
		"projects":   results,
		"categories": categories,
		"count":      cnt,
		"pagination": pagination,
	})
}
