package socket

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"versoul/4lance/auth"
	"versoul/4lance/config"
)

var (
	server *socketio.Server
	conf   = config.GetInstance()
	a      = auth.GetInstance()
)

func init() {
	var err error
	server, err = socketio.NewServer(nil)
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
}

func SocketRoutes(r chi.Router) {
	r.Get("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		handler := server
		handler.ServeHTTP(w, r)
	})
	r.Post("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		handler := server
		handler.ServeHTTP(w, r)
	})
}

func connectHandler(so socketio.Socket, msg string) string {
	dta := map[string]string{}
	err := json.Unmarshal([]byte(msg), &dta)
	if err != nil {
		panic(err)
	}

	so.Join(so.Id())
	mgoSession, err := mgo.Dial(conf.DbHost)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	db := mgoSession.DB(conf.DbName).C("users")

	userData, authOk := a.CheckAuth(dta["sid"])
	if authOk {
		change := bson.M{"$push": bson.M{
			"wids": so.Id(),
		},
		}
		err = db.UpdateId(userData["_id"], change)
		if err != nil {
			panic(err)
		}
	} else {
		query := bson.M{"email": "guest"}
		change := bson.M{"$push": bson.M{
			"wids": so.Id(),
		},
		}
		err = db.Update(query, change)
		if err != nil {
			panic(err)
		}
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

func SendMessageByWid(wid string, project map[string]interface{}) {
	server.BroadcastTo(wid, "newProject", project,
		func(so socketio.Socket, data string) {
			//Клиент подтвердил получение сообщения
			if data == "ok" {
				//fmt.Println("DELIVERED OK")
			}
		})
}
