/*
publisher.go

Creates and publishes digests
*/

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

type Publisher struct {
	reviewFetcher   ReviewFetcher
	outputPath      string
	publishInterval time.Duration
}

// Creates a new publisher based on the given ConfigManager.
// Outputs digests in a folder based on the app ID, thereby allowing for multiple apps
func NewPublisher(configManager ConfigManager) Publisher {
	publisher := Publisher{}
	publisher.outputPath = "./output/" + configManager.AppID()
	publisher.reviewFetcher = NewRSSReviewFetcher(configManager)
	publisher.publishInterval = configManager.PublishInterval()
	return publisher
}

// Returns the file path to the latest published digest if it is within the publish interval
// If there are no recent digests, creates a new one and returns its file path
func (p Publisher) PublishLatest(atTime time.Time) string {
	lastInterval := atTime.Add(-1 * p.publishInterval)
	lastDigestTime, lastDigestPath := p.lastDigest()
	if lastDigestTime.After(lastInterval) {
		return lastDigestPath
	} else {
		return p.writeDigest(lastInterval, atTime)
	}
}

// Fetches reviews from App Store Connect and creates a digest for reviews in a given time window
func (p Publisher) writeDigest(since time.Time, until time.Time) string {
	fmt.Println("Fetching new reviews...")
	reviewsResponse := p.reviewFetcher.FetchReviews(since, until)
	digest := NewDigest(reviewsResponse, since, until)

	publishTimestamp := strconv.FormatInt(until.Unix(), 10)
	fileName := publishTimestamp + ".json"
	filePath := filepath.Join(p.outputPath, fileName)

	content, err := json.Marshal(digest)
	if err != nil {
		log.Fatalf("Error marshalling review content: %v", err)
	}

	os.Remove(filePath)
	fmt.Printf("Creating json at %s\n", filePath)
	err = ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Fatalf("Error writing new digest to %s.\nError: %s", filePath, err)
	}

	// Write the processed markdown file
	mdName := publishTimestamp + ".md"
	mdPath := filepath.Join(p.outputPath, mdName)
	fmt.Printf("Creating markdown at %s\n", mdPath)
	md, err := os.Create(mdPath)
	if err != nil {
		log.Fatalf("Error creating processed markdown file: %s", err)
	}
	defer md.Close()
	_, err = md.WriteString(digest.toMarkdown())
	if err != nil {
		log.Fatalf("Error writing to processed markdown file: %s", err)
	}
	err = md.Sync()
	if err != nil {
		log.Fatalf("Error saving markdown file: %s", err)
	}

	return filePath
}

// Fetches the last digest, if any
// Returns the time it was published and the file path
func (p Publisher) lastDigest() (time.Time, string) {
	// Create the output directory if it doesn't exist
	if _, err := os.Stat(p.outputPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(p.outputPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error making directory %s: %v", p.outputPath, err)
		}
	}

	files, err := ioutil.ReadDir(p.outputPath)
	if err != nil {
		log.Fatalf("Error reading directory %s: %v", p.outputPath, err)
	}

	var latestTimestamp time.Time
	var latestFile string
	// Traverse the directory to find the latest digest
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
			latestFile = filepath.Join(p.outputPath, file.Name())
		}
	}
	return latestTimestamp, latestFile
}
