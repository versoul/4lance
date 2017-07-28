package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"regexp"
	"strings"
	"time"
	"versoul/4lance/delivery"
)

type flParser struct {
	site           string
	dbHost         string
	dbName         string
	collectionName string
}

var fl *flParser

func init() {
	fl = &flParser{
		site:           "https://www.fl.ru",
		dbHost:         "localhost",
		dbName:         "4lance",
		collectionName: "projects",
	}
}
func (self *flParser) parse() {
	session, err := mgo.Dial(self.dbHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(self.dbName).C(self.collectionName)
	doc, err := goquery.NewDocument(self.site + "/projects/?kind=1")
	if err != nil {
		fmt.Println("Ошибка загрузки страницы!")
		fmt.Println(err)
		return
	}
	doc.Find("#projects-list div.b-post:not(.topprjpay)").Each(func(i int, s *goquery.Selection) {

		priceStr := s.Children().First().Text()
		link := s.Find("a.b-post__link")
		projectTitle, err := link.Html()
		checkErr(err)
		projectHref, _ := link.Attr("href")

		re := regexp.MustCompile(`projects\/(.+)\/`)
		projectId := re.FindStringSubmatch(projectHref)[1]

		var projectPrice string
		if strings.Index(priceStr, "По договоренности") == -1 {
			re := regexp.MustCompile(`>(.+)<`)
			res := re.FindStringSubmatch(priceStr)
			projectPrice = strings.TrimSpace(res[1])
		}

		//Find doc with Id
		cnt, err := c.Find(bson.M{"projectId": projectId}).Count()
		if err != nil {
			panic(err)
		}

		//If not exist
		if cnt == 0 {
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
		}
	})
}
func (self *flParser) parseOne(url string, id string) {

	doc, err := goquery.NewDocument(self.site + url)
	if err != nil {
		fmt.Println("Ошибка загрузки страницы!")
		fmt.Println(err)
		return
	}
	descrElem := doc.Find("#projectp" + id)
	projectDescription, err := descrElem.Html()
	checkErr(err)
	projectDescription = strings.TrimSpace(projectDescription)

	catElem := descrElem.Next().Next()
	platnii, err := catElem.Next().Html()
	dateString := ""
	if strings.Index(platnii, "Платный проект") != -1 {
		dateString, err = catElem.Next().Next().Children().Last().Html()
		checkErr(err)
	} else {
		dateString, err = catElem.Next().Children().Last().Html()
		checkErr(err)
	}

	re := regexp.MustCompile(`\[(.+)\]`)
	if re.MatchString(dateString) {
		dateString = re.ReplaceAllString(dateString, "")
	}
	dateString = strings.TrimSpace(dateString)

	loc, _ := time.LoadLocation("Europe/Moscow")
	layout := "02.01.2006 | 15:04"
	projectDate, err := time.ParseInLocation(layout, dateString, loc)
	if err != nil {
		fmt.Println("Project: " + url)
		panic(err)
	}

	var projectCategories []string
	catElem.ChildrenFiltered("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		re := regexp.MustCompile(`s\/(.+)\/`)
		res := re.FindStringSubmatch(href)
		category := strings.TrimSpace(res[1])
		projectCategories = append(projectCategories, category)
	})

	session, err := mgo.Dial(self.dbHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(self.dbName).C(self.collectionName)
	query := bson.M{"projectId": id}
	change := bson.M{"$set": bson.M{"projectDate": projectDate,
		"projectDescription": projectDescription,
		"projectCategories":  projectCategories,
	},
	}
	err = c.Update(query, change)
	if err != nil {
		panic(err)
	}

	delivery.Deliver(id)
}
