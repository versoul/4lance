package auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/gomail.v2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"sync"
	"time"
)

type singleton struct {
	provider string
	config   map[string]interface{}
}

var instance *singleton
var once sync.Once

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}
func (self *singleton) Configure(providerName string, config map[string]interface{}) {
	self.provider = providerName
	self.config = config
}
func (self *singleton) RegisterUser(email string, pass string) (interface{}, error) {
	var id interface{}

	email = strings.ToLower(email)

	switch self.provider {
	case "mongodb":
		mgoSession, err := mgo.Dial(self.config["dbHost"].(string))
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		db := mgoSession.DB(self.config["dbName"].(string)).C("users")
		cnt, err := db.Find(bson.M{"email": email}).Count()
		if err != nil {
			return nil, err

		} else if cnt != 0 {
			return nil, errors.New("Такой email уже зарагестрирован")
		}
		id = bson.NewObjectId()
		err = db.Insert(bson.M{
			"_id":          id,
			"email":        email,
			"pass":         getMD5Hash(pass),
			"confirmEmail": false,
			"admin":        false,
			"wids":         []string{},
		})
		if err != nil {
			panic(err)
		}
	}
	return id, nil
}

func (self *singleton) SendConfirmationEmail(id interface{}) error {

	switch self.provider {
	case "mongodb":
		realId := id.(bson.ObjectId)
		mgoSession, err := mgo.Dial(self.config["dbHost"].(string))
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		db := mgoSession.DB(self.config["dbName"].(string)).C("users")
		var result map[string]interface{}
		err = db.Find(bson.M{"_id": realId}).One(&result)
		if err != nil {
			return err

		}

		email := result["email"].(string)
		confirmationHash := getMD5Hash("4lance.ru" + email)

		query := bson.M{"_id": realId}
		change := bson.M{"$set": bson.M{
			"confirmationHash": confirmationHash,
		},
		}
		err = db.Update(query, change)
		if err != nil {
			panic(err)
		}

		m := gomail.NewMessage()
		m.SetHeader("From", "4lance.ru <"+self.config["mailUser"].(string)+">")
		m.SetHeader("To", email)
		m.SetHeader("Subject", "Подтвердите ваш email")
		m.SetBody("text/html", "<div style='padding:20px 0; font-size:15px;'>"+
			"<p>Добро пожаловать на 4lance.ru</p><p>Для подтверждения email и активации аккаунта, перейдите по ссылке<br/>"+
			"<a href='http://4lance.ru/confirmEmail/"+confirmationHash+"/'>http://4lance.ru/confirmEmail/"+confirmationHash+"/</a></p></div>")

		d := gomail.NewDialer(self.config["mailHost"].(string), self.config["mailPort"].(int),
			self.config["mailUser"].(string), self.config["mailPass"].(string))
		if err := d.DialAndSend(m); err != nil {
			fmt.Println("Ошибка отправки на адрес " + email + " : " + err.Error())
		} else {
			//fmt.Println("Успешно отправлено на адрес " + email)
		}
	}
	return nil
}

func (self *singleton) ConfirmEmail(confirmationHash string) error {

	switch self.provider {
	case "mongodb":
		mgoSession, err := mgo.Dial(self.config["dbHost"].(string))
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		db := mgoSession.DB(self.config["dbName"].(string)).C("users")

		cnt, err := db.Find(bson.M{"confirmationHash": confirmationHash}).Count()
		if err != nil {
			return err
		} else if cnt == 0 {
			return errors.New("Email уже подтвержден")
		}

		query := bson.M{"confirmationHash": confirmationHash}
		change := bson.M{"$set": bson.M{
			"confirmEmail": true,
		}, "$unset": bson.M{
			"confirmationHash": "",
		},
		}
		err = db.Update(query, change)
		if err != nil {
			return err
		}

	}
	return nil
}

func (self *singleton) LoginUser(w http.ResponseWriter, email string, pass string) (map[string]interface{}, error) {
	var result map[string]interface{}

	email = strings.ToLower(email)

	switch self.provider {
	case "mongodb":
		mgoSession, err := mgo.Dial(self.config["dbHost"].(string))
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		db := mgoSession.DB(self.config["dbName"].(string)).C("users")
		cnt, err := db.Find(bson.M{"email": email}).Count()
		if err != nil {
			return nil, err

		} else if cnt == 0 {
			return nil, errors.New("Такой email не зарагестрирован")
		}

		err = db.Find(bson.M{"email": email}).One(&result)
		if err != nil {
			return nil, err

		}

		if getMD5Hash(pass) != result["pass"].(string) {
			return nil, errors.New("Пароль не верный")
		}

		sid, ok := result["sid"].(string)
		if ok {

		} else {
			sidNew, err := uuid.NewV4()
			if err != nil {
				panic(err)
			}
			change := bson.M{"$set": bson.M{
				"sid": sidNew.String(),
			},
			}
			err = db.UpdateId(result["_id"], change)
			sid = sidNew.String()
		}

		expire := time.Now().AddDate(0, 0, 30)
		cookie := http.Cookie{"sid", sid, "/", "",
			expire, "", (86400 * 30), false, false, "", nil}
		http.SetCookie(w, &cookie)
	}
	delete(result, "pass")
	return result, nil
}

func (self *singleton) LogoutUser(w http.ResponseWriter) {
	expire := time.Now().AddDate(-1, 0, 0)
	cookie := http.Cookie{"sid", "", "/", "",
		expire, "", -1, false, false, "", nil}
	http.SetCookie(w, &cookie)
}

func (self *singleton) CheckAuthReq(r *http.Request) (map[string]interface{}, bool) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return nil, false
	}

	return self.CheckAuth(cookie.Value)
}
func (self *singleton) CheckAuth(sid string) (map[string]interface{}, bool) {
	var result map[string]interface{}

	switch self.provider {
	case "mongodb":
		mgoSession, err := mgo.Dial(self.config["dbHost"].(string))
		if err != nil {
			panic(err)
		}
		defer mgoSession.Close()
		db := mgoSession.DB(self.config["dbName"].(string)).C("users")
		/*cnt, err := db.Find(bson.M{"sid": sid}).Count()
		if err != nil {
			return nil, err

		} else if cnt == 0 {
			return nil, errors.New("Такой email не зарагестрирован")
		}*/

		err = db.Find(bson.M{"sid": sid}).One(&result)
		if err != nil {
			return nil, false

		}
	}
	delete(result, "pass")
	return result, true
}
