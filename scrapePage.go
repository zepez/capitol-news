package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type article struct {
	Id         string    `json:"id"`
	Slug       string    `json:"slug"`
	Headline   string    `json:"headline"`
	Subhead    string    `json:"subhead"`
	Body       string    `json:"body"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

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

		// create new from article
		newArticle := article{}

		// create goquery selection
		goquerySelection := e.DOM

		// get date string
		date := goquerySelection.Find("time").Text()

		// remove whitespace and the day of week
		t := CleanDate(date)

		// generate uuid
		newArticle.Id = uuid.Must(uuid.NewRandom()).String()

		// get headline
		newArticle.Headline = goquerySelection.Find("h1").Text()

		// generate slug from headline
		// with random number to ensure unique
		newArticle.Slug = slug.Make(newArticle.Headline) + "-" + strconv.Itoa(rand.Intn(10000))

		// get subhead (sometimes null)
		newArticle.Subhead = goquerySelection.Find("h2").Text()

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
		newArticle.Body = strings.Join(bodyArr, " ")

		// check if article came out today
		if IsToday(t, time.Now()) {
			// check if headline exists
			// if no headline — we don't want it
			if len(newArticle.Headline) > 0 {
				// current timestamp
				newArticle.Created_at = time.Now()
				newArticle.Updated_at = time.Now()

				// generate json from struct
				json, _ := json.Marshal(newArticle)
				buf := bytes.NewBuffer(json)
				fmt.Println(buf)
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
