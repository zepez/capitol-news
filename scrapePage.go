package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gosimple/slug"
)

func ScrapePage() {
	// set true to stop scraping
	stop := false
	sec, _ := strconv.Atoi(os.Getenv("seconds"))
	fmt.Println()

	// start message
	fmt.Println("Scraping for " + os.Getenv("seconds") + " seconds")

	// after defined about of time, stop scraping.
	timer := time.AfterFunc((time.Duration(sec) * time.Second), func() {
		fmt.Println("Stopping...")
		fmt.Println("Done")
		stop = true
	})
	defer timer.Stop()

	// start colly
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

		// generate slug from headline
		// with random number to ensure unique
		slug := slug.Make(headline) + "-" + strconv.Itoa(rand.Intn(10000))

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
		if IsToday(t, time.Now()) {
			// print out
			// check if headline exists
			// if no headline — we don't want it
			if len(headline) > 0 {
				fmt.Println("headline: ", headline)
				fmt.Println("slug: ", slug)
				fmt.Println("subhead: ", subhead)
				fmt.Println("body: ", body)
				fmt.Println("created_at: ", time.Now())
				fmt.Println("updated_at: ", time.Now())
			}
		}
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if !stop {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("============================================================")
		fmt.Println(r.URL)
	})

	c.Visit("https://capitolnewsillinois.com/NEWS/")
}
