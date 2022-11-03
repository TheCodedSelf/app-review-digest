/*
review_fetcher.go

Fetches reviews for an app from App Store Connect
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// This interface abstracts the review fetcher
// In order to improve testability and allow swapping out the fetcher
// i.e. to replace the RSS feed with the newer App Store Connect API
type ReviewFetcher interface {
	FetchReviews(sinceTime time.Time, atTime time.Time) []ReviewsResponseEntry
}

type RSSReviewFetcher struct {
	configManager ConfigManager
}

// Constructs a review fetcher with the given configuration manager
func NewRSSReviewFetcher(configManager ConfigManager) RSSReviewFetcher {
	reviewFetcher := RSSReviewFetcher{}
	reviewFetcher.configManager = configManager
	return reviewFetcher
}

// Fetch reviews from the app's RSS feed.
// Note that this does not implement pagination. This should work fine for an interval of 1 day,
// but fetching longer intervals could miss reviews.
func (r RSSReviewFetcher) FetchReviews(sinceTime time.Time, atTime time.Time) []ReviewsResponseEntry {
	appID := r.configManager.AppID()
	url := fmt.Sprintf("https://itunes.apple.com/us/rss/customerreviews/id=%s/sortBy=mostRecent/page=1/json", appID)
	response, error := http.Get(url)
	if error != nil {
		log.Fatalln(error)
	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatal(error)
	}

	var reviewsResponse ReviewsResponse
	err := json.Unmarshal(body, &reviewsResponse)
	if err != nil {
		log.Fatal(err)
	}

	allReviews := reviewsResponse.Feed.Entry
	fmt.Printf("Fetched %d reviews\n", len(allReviews))
	reviewsForDigest := []ReviewsResponseEntry{}

	for _, review := range allReviews {
		reviewTime, err := time.Parse(time.RFC3339, review.Updated.Label)
		if err != nil {
			log.Fatal(err)
		}
		if reviewTime.After(sinceTime) && reviewTime.Before(atTime) {
			// Since reviews are already ordered, a more efficient approach
			// would be be to find the first review outside the desired range
			// and drop all reviews after that
			reviewsForDigest = append(reviewsForDigest, review)
		}
	}

	return reviewsForDigest
}

type ReviewsResponse struct {
	Feed ReviewsResponseFeed
}

type ReviewsResponseFeed struct {
	Entry []ReviewsResponseEntry
}

type ReviewsResponseEntry struct {
	Title   TextResponse
	Content TextResponse
	Updated TextResponse
	Rating  TextResponse `json:"im:rating"`
}

type TextResponse struct {
	Label string
}
