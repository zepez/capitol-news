package main

import (
	"os"

	"github.com/robfig/cron/v3"
)

func main() {
	os.Setenv("cron", "* * * * *")
	os.Setenv("seconds", "10")

	if len(os.Getenv("cron")) > 0 {
		cj := cron.New()
		cj.AddFunc(os.Getenv("cron"), func() {
			ScrapePage()
		})

		cj.Run()

	}
}
