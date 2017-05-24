package main

import (
	"versoul/4lance/parser"
	"versoul/4lance/routes"
)

var (
	prsr = parser.GetInstance()
)

func main() {

	go prsr.Run()
	routes.InitRoutes()
}
