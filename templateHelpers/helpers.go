package templateHelpers

import (
	"html/template"
	"time"
)

//type helpers struct{}
var Helpers = template.FuncMap{
	"siteToIcon": siteToIcon,
	"siteToName": siteToName,
	"toFullLink": toFullLink,
	"formatTime": formatTime,
	"toHtml":     toHtml,
	"inc": func(i int) int {
		return i + 1
	},
}

func siteToIcon(site string) string {
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
func siteToName(site string) string {
	icon := ""
	switch site {
	case "f-l":
		icon = "fl.ru"
	case "wl":
		icon = "weblancer.net"
	case "fl":
		icon = "freelance.ru"
	case "flm":
		icon = "freelansim.ru"
	}
	return icon
}
func toFullLink(site string, href string) string {
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
func toHtml(str string) template.HTML {
	return template.HTML(str)
}
func formatTime(date time.Time, layout string) string {
	return date.Format(layout)
}
