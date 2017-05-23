package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"html/template"
	"net/http"
	"versoul/4lance/parser"
)

var (
	prsr = parser.GetInstance()
)

func main() {

	go prsr.Run()

	r := gin.Default()
	r.Static("/css", "./static/css")
	r.Static("/lib", "./static/lib")
	r.Static("/js", "./static/js")
	r.Static("/img", "./static/img")
	r.Static("/media", "./static/media")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/google173377f79f6a476a.html", "./static/google173377f79f6a476a.html")
	r.StaticFile("/yandex_87613e29f1d00477.html", "./static/yandex_87613e29f1d00477.html")

	/*r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})*/
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
		err = db.Find(nil).All(&results)
		if err != nil {
			panic(err)
		}

		for _, v := range results {
			if str, ok := v["projectPrice"].(string); ok {
				v["projectPrice"] = template.HTML(str)
			} else {
				//LOL ITS NOT STRING!!
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
