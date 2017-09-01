package templateHelpers

import (
	"html/template"
	"regexp"
	"strings"
	"time"
)

//type helpers struct{}
var Helpers = template.FuncMap{
	"siteToIcon": SiteToIcon,
	"siteToName": SiteToName,
	"toFullLink": ToFullLink,
	"formatTime": FormatTime,
	"toHtml":     ToHtml,
	"stripTags":  StripTags,
	"inc": func(i int) int {
		return i + 1
	},
}

func SiteToIcon(site string) string {
	icon := ""
	switch site {
	case "f-l":
		icon = "free-lance.ru.gif"
	case "wl":
		icon = "weblancer.net.gif"
	case "frl":
		icon = "freelance.png"
	case "flm":
		icon = "freelancim.png"
	}
	return icon
}
func SiteToName(site string) string {
	icon := ""
	switch site {
	case "f-l":
		icon = "fl.ru"
	case "wl":
		icon = "weblancer.net"
	case "frl":
		icon = "freelance.ru"
	case "flm":
		icon = "freelansim.ru"
	}
	return icon
}
func ToFullLink(site string, href string) string {
	link := ""
	switch site {
	case "f-l":
		link = "https://www.fl.ru"
	case "wl":
		link = "https://www.weblancer.net"
	case "frl":
		link = "https://freelance.ru"
	case "flm":
		link = "freelancim.png"
	}
	link += href
	return link
}
func ToHtml(str string) template.HTML {
	return template.HTML(str)
}
func StripTags(str string) string {
	re := regexp.MustCompile(`(?m)<\/?[^>]*>`)
	str = re.ReplaceAllString(str, "")
	str = strings.Replace(str, "&nbsp;", " ", -1)
	return str
}
func FormatTime(date time.Time, layout string) string {
	return date.Format(layout)
}
