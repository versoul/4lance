package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type wlParser struct {
	site           string
	dbHost         string
	dbName         string
	collectionName string
}

var wl *wlParser

func init() {
	wl = &wlParser{
		site:           "https://www.weblancer.net",
		dbHost:         "localhost",
		dbName:         "4lance",
		collectionName: "projects",
	}
}
func (self *wlParser) parse() {
	session, err := mgo.Dial(self.dbHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(self.dbName).C(self.collectionName)
	resp, err := http.Get(self.site + "/jobs/?type=project")
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
	//div.row:not(:has(div>span.featured))
	doc.Find(".cols_table div.row:not(:has(div>span.featured))").Each(func(i int, s *goquery.Selection) {
		projectPrice := strings.TrimSpace(s.Find("div.amount").Text())
		linkElem := s.Find("h2 a")
		projectTitle, err := linkElem.Html()
		projectHref, _ := linkElem.Attr("href")
		checkErr(err)
		dateString, _ := s.Find("div.text-muted span span").Attr("title")
		dateString = strings.TrimSpace(dateString)
		loc, _ := time.LoadLocation("Europe/Kiev")
		layout := "02.01.2006 в 15:04"
		projectDate, err := time.ParseInLocation(layout, dateString, loc)
		checkErr(err)

		re := regexp.MustCompile(`projects\/(.+)\/(.+)\/`)
		res := re.FindStringSubmatch(projectHref)
		projectCategory := res[1]
		projectId := res[2]
		projectCategories := []string{projectCategory}

		//Find doc with Id
		cnt, err := c.Find(bson.M{"projectId": projectId}).Count()
		if err != nil {
			panic(err)
		}

		//If not exist
		if cnt == 0 {
			err = c.Insert(map[string]interface{}{
				"site":              "wl",
				"projectId":         projectId,
				"projectTitle":      projectTitle,
				"projectHref":       projectHref,
				"projectPrice":      projectPrice,
				"projectDate":       projectDate,
				"projectCategories": projectCategories,
			})
			if err != nil {
				log.Fatal(err)
			}
			self.parseOne(projectHref, projectId)
			return
		}
		go self.parseOne(projectHref, projectId)
	})
}
func (self *wlParser) parseOne(url string, id string) {
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
	ii := 1
	projectDescription := ""
	doc.Find(".cols_table .row .col-sm-12").Each(func(i int, s *goquery.Selection) {
		if ii == 2 {
			data, _ := s.Html()
			projectDescription = strings.TrimSpace(data[strings.Index(data, "</div>")+6 : len(data)])
		}
		ii++
	})

	query := bson.M{"projectId": id}
	change := bson.M{"$set": bson.M{
		"projectDescription": projectDescription,
	},
	}
	err = c.Update(query, change)
	if err != nil {
		panic(err)
	}
}
