package main

import (
	"fmt"
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

type Digest struct {
	Title    string
	Subtitle string
	Reviews  []Review
}

type Review struct {
	Title   string
	Content string
	Date    string
	Rating  string
}
