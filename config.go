/*
config.go

Uses a local config.json file to store configuration data
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type TimeOfDay struct {
	Hour   int
	Minute int
}

// Data to persist in config.json
type Configuration struct {
	PublishTime     TimeOfDay
	AppID           string
	PublishInterval time.Duration
}

// If moving to a better solution than a local file, this interface
// allows swapping out the LocalFileConfigManager for i.e. a database instance
type ConfigManager interface {
	PublishTime() TimeOfDay
	SetPublishTime(time TimeOfDay)
	AppID() string
	SetAppID(appID string)
	PublishInterval() time.Duration
	SetPublishInterval(interval time.Duration)
}

// A ConfigManager that stores configuration values in a file.
// Note that this is not atomic; avoid concurrent writes.
type LocalFileConfigManager struct {
	// Where to store the file
	FilePath string
	// If not configured, what app ID to use
	DefaultAppID string
	// If not configured, the length of the digest to be published
	DefaultPublishInterval time.Duration
}

// Creates a LocalFileConfigManager with sensible defaults
func NewLocalFileConfigManager() LocalFileConfigManager {
	return LocalFileConfigManager{
		FilePath:               "./config.json",
		DefaultAppID:           "595068606",
		DefaultPublishInterval: 24 * time.Hour,
	}
}

func (c LocalFileConfigManager) PublishTime() TimeOfDay {
	config := c.fetchConfiguration()
	return config.PublishTime
}

func (c LocalFileConfigManager) SetPublishTime(time TimeOfDay) {
	if time.Hour > 23 {
		log.Fatal(errors.New("Invalid hour for publish time."))
	}

	if time.Minute > 59 {
		log.Fatal(errors.New("Invalid minute for publish time."))
	}

	config := c.fetchConfiguration()
	config.PublishTime = time
	c.writeToFile(config)
}

func (c LocalFileConfigManager) AppID() string {
	config := c.fetchConfiguration()
	if config.AppID == "" {
		c.SetAppID(c.DefaultAppID)
		return c.DefaultAppID
	} else {
		return config.AppID
	}
}

func (c LocalFileConfigManager) SetAppID(appID string) {
	config := c.fetchConfiguration()
	config.AppID = appID
	c.writeToFile(config)
}

func (c LocalFileConfigManager) PublishInterval() time.Duration {
	config := c.fetchConfiguration()
	if config.PublishInterval <= 0 {
		c.SetPublishInterval(c.DefaultPublishInterval)
		return c.DefaultPublishInterval
	} else {
		return config.PublishInterval
	}
}

func (c LocalFileConfigManager) SetPublishInterval(interval time.Duration) {
	config := c.fetchConfiguration()
	config.PublishInterval = interval
	c.writeToFile(config)
}

// Creates the local file if it doesn't exist.
// Returns a Configuration object fetched from the local configuration file.
func (c LocalFileConfigManager) fetchConfiguration() Configuration {
	_, err := os.Stat(c.FilePath)
	if err != nil {
		fmt.Println("Config file doesn't exist yet. Creating now.")
		config := Configuration{}
		config.AppID = c.DefaultAppID
		config.PublishTime = TimeOfDay{}
		config.PublishInterval = c.DefaultPublishInterval
		c.writeToFile(config)
		return config
	}

	// Fetch config from file
	configFile, openError := os.Open(c.FilePath)
	if openError != nil {
		log.Fatal(openError)
	}
	defer configFile.Close()
	byteValue, readError := ioutil.ReadAll(configFile)
	if readError != nil {
		log.Fatal(readError)
	}
	var configuration Configuration
	decodeError := json.Unmarshal(byteValue, &configuration)
	if decodeError != nil {
		log.Fatal(decodeError)
	}
	return configuration
}

// Persist a Configuration object
func (c LocalFileConfigManager) writeToFile(config Configuration) {
	fmt.Print("Writing new configuration to file:")
	fmt.Println(config)

	content, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(c.FilePath)

	err = ioutil.WriteFile(c.FilePath, content, 0644)
	if err != nil {
		log.Fatalf("Error writing config to file %s: %s", c.FilePath, err)
	}
}
