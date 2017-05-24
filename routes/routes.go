package routes

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"html/template"
	"net/http"
	"time"
	"versoul/4lance/templateHelpers"
)

var (
	hlprs = templateHelpers.GetHelpers()
)

func InitRoutes() {
	r := gin.Default()

	initStaticRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	})

	r.LoadHTMLFiles("./templates/base.html",
		"./templates/header.html",
		"./templates/main.html")
	r.GET("/dashboard", func(c *gin.Context) {

		session, err := mgo.Dial("localhost")
		if err != nil {
			panic(err)
		}
		defer session.Close()

		db := session.DB("4lance").C("projects")
		// Query All
		var results []map[string]interface{}
		err = db.Find(nil).Sort("-projectDate").All(&results)
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

		c.HTML(http.StatusOK, "base", gin.H{
			"Error":    "Main website",
			"projects": results,
			"count":    cnt,
		})
	})
	//go r.Run() // listen and serve on 0.0.0.0:8080
	r.Run()
}

func initStaticRoutes(r *gin.Engine) {
	r.Static("/css", "./static/css")
	r.Static("/lib", "./static/lib")
	r.Static("/js", "./static/js")
	r.Static("/img", "./static/img")
	r.Static("/media", "./static/media")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/google173377f79f6a476a.html", "./static/google173377f79f6a476a.html")
	r.StaticFile("/yandex_87613e29f1d00477.html", "./static/yandex_87613e29f1d00477.html")
}
