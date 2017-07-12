package main

import (
	//"github.com/astaxie/beego/session"
	"versoul/4lance/parser"
	"versoul/4lance/routes"
)

var (
	prsr = parser.GetInstance()
	//globalSessions *session.Manager
)

func main() {

	go prsr.Run()
	routes.InitRoutes()
}
