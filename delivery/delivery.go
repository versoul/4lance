package delivery

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"versoul/4lance/config"
	"versoul/4lance/socket"
)

var (
	conf = config.GetInstance()
	sock = socket.GetInstance()
)

func Deliver(projectId string) {
	project := getProjectById(projectId)

	mgoSession, err := mgo.Dial(conf.DbHost)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	dbC := mgoSession.DB(conf.DbName).C("users")
	// Query All
	var users []map[string]interface{}
	query := bson.M{"$and": []bson.M{
		bson.M{"$where": "this.wids.length>0"},
		bson.M{"filter.categories": bson.M{"$in": project["projectCategories"]}},
	},
	}
	err = dbC.Find(query).Select(bson.M{"_id": 0}).All(&users)
	if err != nil {
		panic(err)
	}

	wids := []interface{}{}
	for _, user := range users {
		wids = user["wids"].([]interface{})
		for _, wid := range wids {
			fmt.Println("Send to socket wid=" + wid.(string))
			//deliverBySocket(wid.(string), project)
			sock.SendMessageByWid(wid.(string), project)
		}
	}
}

func getProjectById(projectId string) map[string]interface{} {
	mgoSession, err := mgo.Dial(conf.DbHost)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	dbC := mgoSession.DB(conf.DbName).C("projects")
	// Query All
	var project map[string]interface{}
	query := bson.M{"projectId": projectId}
	err = dbC.Find(query).Select(bson.M{"_id": 0}).One(&project)
	if err != nil {
		panic(err)
	}
	return project
}
