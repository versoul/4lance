package parser

import (
	"sync"
	"time"
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

type singleton struct{}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

func (self *singleton) Run() {
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		go fl.parse()
		go wl.parse()
		go frl.parse()
	}
}
