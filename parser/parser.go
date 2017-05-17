package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
	"sync"
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

type singleton struct {
	site string
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{
			site: "https://www.fl.ru",
		}
	})
	return instance
}

func (self *singleton) Parse() {
	doc, err := goquery.NewDocument(self.site + "/projects/?kind=1")
	checkErr(err)
	doc.Find("#projects-list div.b-post:not(.topprjpay)").Each(func(i int, s *goquery.Selection) {
		priceStr := s.Children().First().Text()
		//TODO if not find id only!
		link := s.Find("a.b-post__link")
		projectTitle, err := link.Html()
		checkErr(err)
		projectHref, _ := link.Attr("href")

		fmt.Println(projectTitle + "--" + projectHref)

		var projectPrice string
		if strings.Index(priceStr, "По договоренности") == -1 {
			re := regexp.MustCompile(`>(.+)<`)
			res := re.FindStringSubmatch(priceStr)
			projectPrice = strings.TrimSpace(res[1])
		}
		fmt.Println(projectPrice)

		re := regexp.MustCompile(`projects\/(.+)\/`)
		projectId := re.FindStringSubmatch(projectHref)[1]
		fmt.Println(projectId)

		fmt.Println("***")
	})

}
