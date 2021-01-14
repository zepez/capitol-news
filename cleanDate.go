package main

import (
	"fmt"
	"strings"
	"time"
)

func CleanDate(d string) time.Time {

	// split date string on spaces
	dateArr := strings.SplitAfter(d, " ")
	fmt.Println(dateArr)

	var dateArrC []string
	// clean slice of empties
	for _, str := range dateArr {
		if str != " " {
			dateArrC = append(dateArrC, str)
		}
	}

	// remove first in array ("Monday, ", "Tuesday, " etc)
	dateArrC = append(dateArrC[:0], dateArrC[0+1:]...)
	// join arr back together
	date := strings.Join(dateArrC, "")

	// define time parsing layout
	const (
		layoutUS = "January 2, 2006"
	)

	// parse time
	t, _ := time.Parse(layoutUS, date)

	return t

}
