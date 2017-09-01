package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"golang.org/x/net/html/charset"
	"net/http"
	"regexp"
	"strings"
	"time"
	"versoul/4lance/delivery"
)

type frlParser struct {
	site           string
	dbHost         string
	dbName         string
	collectionName string
}

var frl *frlParser

func init() {
	frl = &frlParser{
		site:           "https://freelance.ru",
		dbHost:         "localhost",
		dbName:         "4lance",
		collectionName: "projects",
	}
}
func (self *frlParser) parse() {
	session, err := mgo.Dial(self.dbHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(self.dbName).C(self.collectionName)
	resp, err := http.Get(self.site + "/projects/")
	if err != nil {
		fmt.Println("Ошибка загрузки страницы!")
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		checkErr(err)
	}
	doc, err := goquery.NewDocumentFromReader(utf8)

	doc.Find(".projects .proj.public:not(.prio)").Each(func(i int, s *goquery.Selection) { //:not(.prio) поднятые
		link := s.Find("a.ptitle")
		projectHref, _ := link.Attr("href")
		projectTitle, _ := link.Find("span").Html()

		re := regexp.MustCompile(`projects\/(.+)\.html`)
		projectId := re.FindStringSubmatch(projectHref)[1]

		priceStr, _ := s.Find(".cost a").Html()
		priceStr = strings.TrimSpace(priceStr)
		var projectPrice string
		if strings.Index(priceStr, "Договорная") == -1 {
			projectPrice = priceStr
		}

		projectDescription, _ := s.Find(".descr p span").Last().Html()

		dateStr, _ := s.Find(".pdata").Attr("title")
		dateStr = strings.TrimSpace(dateStr)
		loc, _ := time.LoadLocation("Europe/Moscow")
		layout := "02.01.06 15:04"
		projectDate, err := time.ParseInLocation(layout, dateStr, loc)
		if err != nil {
			fmt.Println("Project: " + projectHref)
			panic(err)
		}

		//Find doc with Id
		cnt, err := c.Find(bson.M{"projectId": projectId}).Count()
		if err != nil {
			panic(err)
		}

		//If not exist
		if cnt == 0 {
			err = c.Insert(map[string]interface{}{
				"site":               "frl",
				"projectId":          projectId,
				"projectTitle":       projectTitle,
				"projectHref":        projectHref,
				"projectPrice":       projectPrice,
				"projectDate":        projectDate,
				"projectDescription": projectDescription,
			})
			if err != nil {
				panic(err)
			}
			go self.parseOne(projectHref, projectId)
		}
	})

}
func (self *frlParser) parseOne(url string, id string) {

	session, err := mgo.Dial(self.dbHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(self.dbName).C(self.collectionName)
	resp, err := http.Get(self.site + url)
	if err != nil {
		fmt.Println("Ошибка загрузки страницы!")
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		checkErr(err)
	}
	doc, err := goquery.NewDocumentFromReader(utf8)

	projectDescription, _ := doc.Find("#proj_table tr").Next().Find("p").Html()

	categoryHref, _ := doc.Find(".breadcrumb li").Last().Find("a").Attr("href")
	re := regexp.MustCompile(`projects\/\?spec=(.+)`)
	cat := re.FindStringSubmatch(categoryHref)[1]

	var projectCategories []string
	projectCategories = append(projectCategories, cat)

	query := bson.M{"projectId": id}
	change := bson.M{"$set": bson.M{
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
