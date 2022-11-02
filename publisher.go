package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// For demonstration purposes, consider changing the interval period to once per minute:
// const PUBLISH_INTERVAL_PERIOD time.Duration = 1 * time.Minute

// Or, to populate over a greater time frame:
// const PUBLISH_INTERVAL_PERIOD time.Duration = 14 * 24 * time.Hour

// Or, leave at 24 hours.
const PUBLISH_INTERVAL_PERIOD time.Duration = 24 * time.Hour

type Publisher struct {
	ReviewFetcher ReviewFetcher
	OutputPath    string
}

func NewPublisher(configManager ConfigManager) Publisher {
	publisher := Publisher{}
	publisher.OutputPath = "./output/" + configManager.GetAppID()
	publisher.ReviewFetcher = NewRSSReviewFetcher(configManager)
	return publisher
}

// If there is a digest at publish time in the last 24 hours, use that
// If not, make a digest and write it at publish time
func (p Publisher) PublishLatest(atTime time.Time) string {
	lastInterval := atTime.Add(-1 * PUBLISH_INTERVAL_PERIOD)
	lastDigestTime, lastDigestPath := p.lastDigest()
	if lastDigestTime.After(lastInterval) {
		return lastDigestPath
	} else {
		return p.writeDigest(lastInterval, atTime)
	}
}

func (p Publisher) writeDigest(since time.Time, until time.Time) string {
	fmt.Println("Fetching new reviews...")
	reviewsResponse := p.ReviewFetcher.FetchReviews(since, until)
	digest := NewDigest(reviewsResponse, since, until)

	publishTimestamp := strconv.FormatInt(until.Unix(), 10)
	fileName := publishTimestamp + ".json"
	filePath := filepath.Join(p.OutputPath, fileName)

	content, err := json.Marshal(digest)
	if err != nil {
		log.Fatalf("Error marshalling review content: %v", err)
	}

	os.Remove(filePath)

	err = ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Fatalf("Error writing new digest to %s.\nError: %v", filePath, err)
	}
	return filePath
}

func (p Publisher) lastDigest() (time.Time, string) {
	if _, err := os.Stat(p.OutputPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(p.OutputPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error making directory %s: %v", p.OutputPath, err)
		}
	}

	files, err := ioutil.ReadDir(p.OutputPath)
	if err != nil {
		log.Fatalf("Error reading directory %s: %v", p.OutputPath, err)
	}

	var latestTimestamp time.Time
	var latestFile string
	for _, file := range files {
		fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		fileTimestamp, err := strconv.ParseInt(fileName, 10, 64)
		if err != nil {
			fmt.Printf("Error getting timestamp integer from file name.\nFile name: %v\nError: %v\n", fileName, err)
			continue
		}
		digestTime := time.Unix(fileTimestamp, 0)

		if digestTime.After(latestTimestamp) {
			latestTimestamp = digestTime
			latestFile = filepath.Join(p.OutputPath, file.Name())
		}
	}
	return latestTimestamp, latestFile
}
