package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"versoul/4lance/parser"
	//"time"
)

var (
	prsr = parser.GetInstance()
)

func main() {

	prsr.Parse()

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
		"./templates/error.html")
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base", gin.H{
			"Error": "Main website",
		})
	})
	//go r.Run() // listen and serve on 0.0.0.0:8080
	r.Run()

	/*ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		fmt.Println("TIME")
	}*/
}
