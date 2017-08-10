package delivery

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"versoul/4lance/config"
	"versoul/4lance/socket"
	"versoul/4lance/templateHelpers"
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
	var usersFiltered []map[string]interface{}
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
			usersFiltered = append(usersFiltered, user)
		}
	}

	//Проходим по пользователям которым подошел этот проект
	for _, user := range usersFiltered {
		//Если отправка по вебсокету
		wids = user["wids"].([]interface{})
		for _, wid := range wids {
			//fmt.Println("Send wid=" + wid.(string))
			//deliverBySocket(wid.(string), project)
			socket.SendMessageByWid(wid.(string), "newProject", project)
		}
		//если через pushAll
		//TODO еще передавать ид юзера в пушал

		paid, ok := user["paid"].(string)
		if ok && paid != "" {
			sendByPushAll(paid, project)
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

func sendByPushAll(paid string, project map[string]interface{}) {
	link := "https://pushall.ru/api.php"
	//https://pushall.ru/api.php?type=unicast&id=3682&key=e093cedde9bde238d20ebd23bbbd2ac6&uid=59580&title=test&text=testtext\nnewtext&url=http://4lance.ru/dashboard/
	//var jsonStr = `{"type":"unicast", "id":"3682", "key":"e093cedde9bde238d20ebd23bbbd2ac6", "uid":59580, "title":"` + project["projectTitle"].(string) + `", "text":"` + project["projectDescription"].(string) + `", "url": "http://4lance.ru/dashboard/"}`
	fmt.Println(project["projectPrice"].(string))
	fmt.Println(templateHelpers.StripTags(project["projectPrice"].(string)))
	form := url.Values{
		"type":  {"unicast"},
		"id":    {"3682"},
		"key":   {"e093cedde9bde238d20ebd23bbbd2ac6"},
		"uid":   {paid},
		"title": {" " + templateHelpers.StripTags(project["projectTitle"].(string))},
		"text": {
			templateHelpers.StripTags(project["projectDescription"].(string)) +
				" Бюджет: " + templateHelpers.StripTags(project["projectPrice"].(string)),
		},
		"url": {templateHelpers.ToFullLink(project["site"].(string), project["projectHref"].(string))},
	}
	/*form.Add("type", "unicast")
	form.Add("id", "3682")
	form.Add("key", "e093cedde9bde238d20ebd23bbbd2ac6")
	form.Add("uid", "59580")
	form.Add("title", templateHelpers.StripTags(project["projectTitle"].(string)))
	form.Add("text", templateHelpers.StripTags(project["projectDescription"].(string)))
	form.Add("url", templateHelpers.ToFullLink(project["site"].(string), project["projectHref"].(string)))*/

	req, err := http.NewRequest("POST", link, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
