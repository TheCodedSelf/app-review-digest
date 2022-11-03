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

func NewPublisher(configManager ConfigManager) Publisher {
	publisher := Publisher{}
	publisher.outputPath = "./output/" + configManager.AppID()
	publisher.reviewFetcher = NewRSSReviewFetcher(configManager)
	publisher.publishInterval = configManager.PublishInterval()
	return publisher
}

// If there is a digest at publish time since the last publish interval, use that
// If not, make a digest and write it at publish time
func (p Publisher) PublishLatest(atTime time.Time) string {
	lastInterval := atTime.Add(-1 * p.publishInterval)
	lastDigestTime, lastDigestPath := p.lastDigest()
	if lastDigestTime.After(lastInterval) {
		return lastDigestPath
	} else {
		return p.writeDigest(lastInterval, atTime)
	}
}

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

func (p Publisher) lastDigest() (time.Time, string) {
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
