package main

import (
	//"github.com/astaxie/beego/session"
	"sync"
	"versoul/4lance/parser"
	"versoul/4lance/routes"
)

var (
	prsr = parser.GetInstance()
	wg   sync.WaitGroup
	//globalSessions *session.Manager
)

func main() {

	wg.Add(1)
	go func() {
		defer wg.Done()
		prsr.Run()
	}()
	routes.InitRoutes()
	wg.Wait()
}
