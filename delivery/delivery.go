package delivery

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"versoul/4lance/config"
	"versoul/4lance/socket"
)

var (
	conf = config.GetInstance()
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
		userFilter := user["filter"].(map[string]interface{})
		userKeywords := userFilter["keywords"].([]interface{})
		findKeyword := false
		for _, keyword := range userKeywords {
			re := regexp.MustCompile(`(?mi)` + keyword.(string))
			if re.MatchString(project["projectDescription"].(string)) {
				findKeyword = true
			}
			if re.MatchString(project["projectTitle"].(string)) {
				findKeyword = true
			}
		}

		if len(userKeywords) == 0 || findKeyword {
			wids = user["wids"].([]interface{})
			for _, wid := range wids {
				//fmt.Println("Send wid=" + wid.(string))
				//deliverBySocket(wid.(string), project)
				socket.SendMessageByWid(wid.(string), project)
			}
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
