package main

import (
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
)

func main() {
	os.Setenv("cron", "* * * * *")
	if len(os.Getenv("cron")) > 0 {
		cj := cron.New()
		cj.AddFunc(os.Getenv("cron"), func() {
			fmt.Println("Starting script")
			ScrapePage()
		})

		cj.Run()
	}
}
