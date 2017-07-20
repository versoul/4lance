package auth

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

type singleton struct {
	provider string
	config   map[string]interface{}
}

var instance *singleton
var once sync.Once

func GetInstance(providerName string, config map[string]interface{}) *singleton {
	once.Do(func() {
		instance = &singleton{
			provider: providerName,
			config:   config,
		}
	})
	return instance
}

func (self *singleton) RegisterUser(email string, pass string) error {
	fmt.Println(self.provider)
	switch self.provider {
	case "mongodb":
		mgoSession, err := mgo.Dial(self.config["host"].(string))
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		//Get data
		db := mgoSession.DB(self.config["dbName"].(string)).C("users")

		//Find doc with Id
		cnt, err := db.Find(bson.M{"email": email}).Count()
		if err != nil && cnt != 0 {
			return errors.New("Email already exists.")
		}
		fmt.Println("COUNT")
		fmt.Println(cnt)
		//If not exist
		/*if cnt == 0 {
			err = c.Insert(map[string]interface{}{
				"site":         "f-l",
				"projectId":    projectId,
				"projectTitle": projectTitle,
				"projectHref":  projectHref,
				"projectPrice": projectPrice,
			})
			if err != nil {
				log.Fatal(err)
			}
			go self.parseOne(projectHref, projectId)
		}*/

	}
	return nil
}
