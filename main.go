package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("LOL")
	r := gin.Default()
	r.Static("/css", "./static/css")
	r.Static("/lib", "./static/lib")
	r.Static("/js", "./static/js")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.LoadHTMLFiles("./templates/base.html",
		"./templates/header.html",
		"./templates/error.html")
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base", gin.H{
			"Error": "Main website",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
