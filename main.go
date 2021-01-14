package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		// specify allowed domains
		colly.AllowedDomains("www.capitolnewsillinois.com", "capitolnewsillinois.com"),
		colly.URLFilters(
			regexp.MustCompile("/NEWS/"),
		),
	)

	// reate limit
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*capitolnewsillinois*",
		Parallelism: 2,
		Delay:       5 * time.Second,
	})

	// get date
	c.OnHTML(".edn_article, .edn_articleDetails", func(e *colly.HTMLElement) {
		// get date string
		date := e.ChildText(".edn_metaDetails")
		// split date string on spaces
		dateArr := strings.SplitAfter(date, " ")
		// remove first in array ("Monday, ", "Tuesday, " etc)
		dateArr = append(dateArr[:0], dateArr[0+1:]...)
		// join arr back together
		date = strings.Join(dateArr, " ")

		const (
			layoutUS = "January 2, 2006"
		)

		t, _ := time.Parse(layoutUS, date)

		// create goquery selection
		goquerySelection := e.DOM
		// get headline
		headline := goquerySelection.Find("h1").Text()
		// get subhead (sometimes null)
		subhead := goquerySelection.Find("h2").Text() // subhead
		// create body slice
		var bodyArr []string
		// get all paragraph tags
		// map through and add to slice
		goquerySelection.Find("p").Each(func(i int, s *goquery.Selection) {
			// gets inner html
			temp, _ := s.Html()
			// add paragraph tags to outer
			temp = "<p>" + temp + "</p>"
			// add to body slice
			bodyArr = append(bodyArr, temp)
		})
		// convert slice to string with spaces
		body := strings.Join(bodyArr, " ")

		// check if article came out today
		if isToday(t, time.Now()) {
			// @todo. nothing came out today, so can't test lol
			// print out
			if len(headline) > 0 {
				fmt.Println("headline: ", headline)
				fmt.Println("subhead: ", subhead)
				fmt.Println("body: ", body)
			}
		}
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
		fmt.Println("============================================================")
	})

	c.Visit("https://capitolnewsillinois.com/NEWS/")
}

func isToday(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
