package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github"
)

// !+
// func main() {
// 	result, err := github.SearchIssues(os.Args[1:])
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("%d issues:\n", result.TotalCount)
// 	for _, item := range result.Items {
// 		fmt.Printf("#%-5d %9.9s %.55s\n",
// 			item.Number, item.User.Login, item.Title)
// 	}
// }

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	now := time.Now()

	var lessThanMonthOld, lessThanYearOld, moreThanYearOld []*github.Issue

	for _, item := range result.Items {
		age := now.Sub(item.CreatedAt)

		switch {
		case age < 31*24*time.Hour:
			lessThanMonthOld = append(lessThanMonthOld, item)
		case age < 365*24*time.Hour:
			lessThanYearOld = append(lessThanYearOld, item)
		default:
			moreThanYearOld = append(moreThanYearOld, item)
		}
	}
	printIssues("Less than a month old", lessThanMonthOld)
	printIssues("Less than a year old", lessThanYearOld)
	printIssues("More than a yaer old", moreThanYearOld)
}

func printIssues(title string, issues []*github.Issue) {
	if len(issues) > 0 {
		fmt.Printf("\n%s:\n", title)
		for _, item := range issues {
			fmt.Printf("#%-5d %9.9s %.55s\n",
				item.Number, item.User.Login, item.Title)
		}
	}
}
