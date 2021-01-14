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
		colly.AllowedDomains("capitolnewsillinois.com"),
		colly.URLFilters(
			regexp.MustCompile("/NEWS/"),
		),
	)

	// rate limit
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*capitolnewsillinois*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	// get date
	c.OnHTML(".edn_article, .edn_articleDetails", func(e *colly.HTMLElement) {

		// create goquery selection
		goquerySelection := e.DOM

		// get date string
		// date := e.ChildText(".edn_metaDetails")
		date := goquerySelection.Find("time").Text()

		// remove whitespace and the day of week
		t := CleanDate(date)

		// get headline
		headline := goquerySelection.Find("h1").Text()
		// get subhead (sometimes null)
		subhead := goquerySelection.Find("h2").Text()
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
			// print out
			fmt.Println("TODAY")
			// check if headline exists
			// if no headline — we don't want it
			if len(headline) > 0 {
				fmt.Println("headline: ", headline)
				fmt.Println("subhead: ", subhead)
				fmt.Println("body: ", body)
				fmt.Println("created_at: ", time.Now())
				fmt.Println("updated_at: ", time.Now())
			}
		}
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("============================================================")
		fmt.Println(r.URL)
	})

	c.Visit("https://capitolnewsillinois.com/NEWS/")
}

// check if today
func isToday(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
