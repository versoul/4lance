package templateHelpers

import (
	"sync"
	"time"
)

type helpers struct{}

var instance *helpers
var once sync.Once

func GetHelpers() *helpers {
	once.Do(func() {
		instance = &helpers{}
	})
	return instance
}
func (self *helpers) SiteToIcon(site string) string {
	icon := ""
	switch site {
	case "f-l":
		icon = "free-lance.ru.gif"
	case "wl":
		icon = "weblancer.net.gif"
	case "fl":
		icon = "freelance.png"
	case "flm":
		icon = "freelancim.png"
	}
	return icon
}
func (self *helpers) ToFullLink(site string, href string) string {
	link := ""
	switch site {
	case "f-l":
		link = "https://www.fl.ru"
	case "wl":
		link = "https://www.weblancer.net"
	case "fl":
		link = "freelance.png"
	case "flm":
		link = "freelancim.png"
	}
	link += href
	return link
}
func (self *helpers) FormatTime(date time.Time, layout string) string {
	return date.Format(layout)
}
