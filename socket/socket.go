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
func connectHandler(so socketio.Socket, msg string) string {
	dta := map[string]string{}
	err := json.Unmarshal([]byte(msg), &dta)
	if err != nil {
		panic(err)
	}

	userData, authOk := a.CheckAuth(dta["sid"])
	if authOk {
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

/*
func sendMessageByWid(args *fromModuleDta) {
	var recordId int
	db.Connect()
	err := db.DB.QueryRow(`INSERT INTO messages (wid, message, mtype)
		VALUES($1,$2, $3) returning id;`,
		&args.To.Wid, &args.MessageString, &args.Mtype).Scan(&recordId)
	utils.CheckErr(err)
	db.Close()
	var clientDta = toClientDta{
		Message: args.MessageString,
		Mtype:   args.Mtype,
	}
	Server.BroadcastTo(args.To.Wid, "msg", clientDta,
		func(so socketio.Socket, data string) {
			//Клиент подтвердил получение сообщения
			if data == "ok" {
				//log.Println("get answer from client")
				db.Connect()
				_, err = db.DB.Query("DELETE FROM messages WHERE id = $1",
					recordId)
				utils.CheckErr(err)
				db.Close()
			}
		})
}
*/
