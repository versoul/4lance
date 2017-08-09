package socket

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"github.com/pressly/chi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
	"versoul/4lance/auth"
	"versoul/4lance/config"
)

var (
	server     *socketio.Server
	conf       = config.GetInstance()
	a          = auth.GetInstance()
	activeWids = []string{}
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

	ticker := time.NewTicker(1 * time.Minute)
	go socketPingLoop(ticker)

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

func SendMessageByWid(wid string, msgType string, data map[string]interface{}) {
	server.BroadcastTo(wid, msgType, data,
		func(so socketio.Socket, data string) {
			//Клиент подтвердил получение сообщения
			if msgType == "pingSocket" && data == "ok" {
				removeFromActiveWids(so.Id())
			}
		})
}

func socketPingLoop(ticker *time.Ticker) {
	for range ticker.C {
		go socketPing()
	}
}

func socketPing() {
	if len(activeWids) > 0 {
		removeWidsFromDb()
		activeWids = []string{}
	}

	mgoSession, err := mgo.Dial(conf.DbHost)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	db := mgoSession.DB(conf.DbName).C("users")
	result := []map[string]interface{}{}
	query := bson.M{"$where": "this.wids.length>0"}
	err = db.Find(query).Select(bson.M{"wids": 1, "email": 1, "_id": -1}).All(&result)
	if err != nil {
		panic(err)
	}

	for _, val := range result {
		usrWids := val["wids"].([]interface{})
		for _, wid := range usrWids {
			activeWids = append(activeWids, wid.(string))
		}
	}

	for _, wid := range activeWids {
		SendMessageByWid(wid, "pingSocket", map[string]interface{}{})
	}
}
func removeFromActiveWids(wid string) {
	for i, val := range activeWids {
		if val == wid {
			activeWids = append(activeWids[:i], activeWids[i+1:]...)
		}
	}
}
func removeWidsFromDb() {
	mgoSession, err := mgo.Dial(conf.DbHost)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	db := mgoSession.DB(conf.DbName).C("users")
	query := bson.M{"wids": bson.M{"$in": activeWids}}
	change := bson.M{"$pull": bson.M{
		"wids": bson.M{"$in": activeWids},
	},
	}
	db.Update(query, change)
}
