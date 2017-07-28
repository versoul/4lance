package socket

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"sync"
	"versoul/4lance/auth"
	"versoul/4lance/config"
)

var (
	conf = config.GetInstance()
	a    = auth.GetInstance()
)

type singleton struct {
	server *socketio.Server
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		server, err := socketio.NewServer(nil)
		if err != nil {
			panic(err)
		}
		server.On("connection", func(so socketio.Socket) {
			//fmt.Println(server.GetMaxConnection())
			//Максимальное кол-воподключений
			//Возможнонужно будет покрутить цифру в +
			// по умолчанию 1000 максимально
			so.On("conn", connectHandler)
			so.On("disconnection", disconnectHandler)
		})
		server.On("error", func(so socketio.Socket, err error) {
			panic(err)
		})
		fmt.Println("SERVER")
		fmt.Println(server)
		instance = &singleton{
			server: server,
		}
	})
	return instance
}

func (self *singleton) SocketRoutes(r chi.Router) {
	r.Get("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		handler := self.server
		handler.ServeHTTP(w, r)
	})
	r.Post("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		handler := self.server
		handler.ServeHTTP(w, r)
	})
}

func connectHandler(so socketio.Socket, msg string) string {
	dta := map[string]string{}
	err := json.Unmarshal([]byte(msg), &dta)
	if err != nil {
		panic(err)
	}

	userData, authOk := a.CheckAuth(dta["sid"])
	if authOk {
		so.Join(so.Id())
		mgoSession, err := mgo.Dial(conf.DbHost)
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		db := mgoSession.DB(conf.DbName).C("users")
		change := bson.M{"$push": bson.M{
			"wids": so.Id(),
		},
		}
		err = db.UpdateId(userData["_id"], change)
		if err != nil {
			panic(err)
		}

		instance.SendMessageByWid(so.Id(), map[string]interface{}{"lol": "lol"})

	}
	return "ok"
}
func disconnectHandler(so socketio.Socket) {
	mgoSession, err := mgo.Dial(conf.DbHost)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	db := mgoSession.DB(conf.DbName).C("users")
	query := bson.M{"wids": so.Id()}
	change := bson.M{"$pull": bson.M{
		"wids": so.Id(),
	},
	}
	db.Update(query, change)
}

func (self *singleton) SendMessageByWid(wid string, project map[string]interface{}) {
	fmt.Println("send " + wid)
	fmt.Println("SERVER")
	fmt.Println(self.server)
	self.server.BroadcastTo(wid, "msg", project,
		func(so socketio.Socket, data string) {
			//Клиент подтвердил получение сообщения
			if data == "ok" {
				fmt.Println("DELIVERED OK")
			}
		})
}
