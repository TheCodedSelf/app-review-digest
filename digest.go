/*
digest.go
Data object for digests.
Includes a constructor and a method to create a markdown representation.
*/
package main

import (
	"fmt"
	"strings"
	"time"
)

type Digest struct {
	Title    string
	Subtitle string
	Reviews  []Review
}

type Review struct {
	Author  string
	Title   string
	Content string
	Date    string
	Rating  string
}

func NewDigest(reviewResponses []ReviewsResponseEntry, since time.Time, until time.Time) Digest {
	// Only includes the year in the 'to' date
	fromString := since.Format("Mon Jan _2 15:04:05")
	toString := until.Format(time.ANSIC)

	title := fmt.Sprintf("App reviews from %s until %s", fromString, toString)
	var subtitle string
	if len(reviewResponses) == 0 {
		subtitle = "No reviews today, try again tomorrow!"
	} else {
		subtitle = fmt.Sprintf("%d reviews", len(reviewResponses))
	}
	var reviews []Review

	// Populate reviews from response objects
	for _, reviewResponse := range reviewResponses {
		review := Review{
			Author:  reviewResponse.Author.Name.Label,
			Title:   reviewResponse.Title.Label,
			Content: reviewResponse.Content.Label,
			Date:    reviewResponse.Updated.Label,
			Rating:  reviewResponse.Rating.Label,
		}
		reviews = append(reviews, review)
	}
	return Digest{
		Title:    title,
		Subtitle: subtitle,
		Reviews:  reviews,
	}
}

// Creates a markdown string representing the digest
func (d Digest) ToMarkdown() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("# %s\n", d.Title))
	builder.WriteString(fmt.Sprintf("%s\n", d.Subtitle))
	for _, review := range d.Reviews {
		builder.WriteString(fmt.Sprintf("## %s\n", review.Title))
		date, err := time.Parse(time.RFC3339, review.Date)
		var dateString string
		if err == nil {
			// A more readable date.
			dateString = date.Format(time.ANSIC)
		} else {
			dateString = review.Date
		}
		builder.WriteString(fmt.Sprintf("**%s star(s)** — _%s_\n\n", review.Rating, dateString))
		builder.WriteString(fmt.Sprintf("%s\n", review.Content))
		builder.WriteString(fmt.Sprintf("\nby _%s_\n", review.Author))
	}
	return builder.String()
}
