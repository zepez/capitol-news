package main

import (
	"os"

	"github.com/robfig/cron/v3"
)

func main() {
	// env for dev
	// os.Setenv("cron", "* * * * *")
	// os.Setenv("seconds", "60")
	// os.Setenv("endpoint", "http://localhost:3001")

	if len(os.Getenv("cron")) > 0 {
		cj := cron.New()
		cj.AddFunc(os.Getenv("cron"), func() {
			ScrapePage()
		})

		cj.Run()
	}
}
