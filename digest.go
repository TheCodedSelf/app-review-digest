package main

import (
	"fmt"
	"strings"
	"time"
)

func NewDigest(reviewResponses []ReviewsResponseEntry, since time.Time, until time.Time) Digest {
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
	for _, reviewResponse := range reviewResponses {
		review := Review{
			title:   reviewResponse.Title.Label,
			content: reviewResponse.Content.Label,
			date:    reviewResponse.Updated.Label,
			rating:  reviewResponse.Rating.Label,
		}
		reviews = append(reviews, review)
	}
	return Digest{
		title:    title,
		subtitle: subtitle,
		reviews:  reviews,
	}
}

type Digest struct {
	title    string
	subtitle string
	reviews  []Review
}

func (d Digest) toMarkdown() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("# %s\n", d.title))
	builder.WriteString(fmt.Sprintf("%s\n", d.subtitle))
	for _, review := range d.reviews {
		builder.WriteString(fmt.Sprintf("## %s\n", review.title))
		date, err := time.Parse(time.RFC3339, review.date)
		var dateString string
		if err == nil {
			// A more readable date.
			dateString = date.Format(time.ANSIC)
		} else {
			dateString = review.date
		}
		builder.WriteString(fmt.Sprintf("**%s star(s)** — _%s_\n\n", review.rating, dateString))
		builder.WriteString(fmt.Sprintf("%s\n", review.content))
	}
	return builder.String()
}

type Review struct {
	title   string
	content string
	date    string
	rating  string
}
