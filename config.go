package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Configuration struct {
	PublishTime TimeOfDay
	AppID       string
}

type TimeOfDay struct {
	Hour   int
	Minute int
}

type ConfigManager interface {
	PublishTime() TimeOfDay
	SetPublishTime(time TimeOfDay)
	AppID() string
	SetAppID(appID string)
}

type LocalFileConfigManager struct {
	FilePath     string
	DefaultAppID string
}

func NewLocalFileConfigManager() LocalFileConfigManager {
	return LocalFileConfigManager{
		FilePath:     "./config.json",
		DefaultAppID: "595068606",
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

func (c LocalFileConfigManager) fetchConfiguration() Configuration {
	_, err := os.Stat(c.FilePath)
	if err != nil {
		fmt.Println("Config file doesn't exist yet. Creating now.")
		config := Configuration{}
		config.AppID = c.DefaultAppID
		config.PublishTime = TimeOfDay{}
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
