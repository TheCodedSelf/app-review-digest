package main

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestSetsAndReturnsPublishTime(t *testing.T) {
	expected := TimeOfDay{
		Hour:   rand.Intn(23),
		Minute: rand.Intn(59),
	}
	configManager := newConfigManager()
	configManager.SetPublishTime(expected)

	actual := configManager.PublishTime()

	if expected != actual {
		t.Fail()
	}
}

func TestSetsAndReturnsAppID(t *testing.T) {
	expected := strconv.FormatInt(rand.Int63(), 10)
	configManager := newConfigManager()
	configManager.SetAppID(expected)

	actual := configManager.AppID()

	if expected != actual {
		t.Fail()
	}
}

func TestSetsAndReturnsPublishInterval(t *testing.T) {
	expected := time.Duration(rand.Int63())
	configManager := newConfigManager()
	configManager.SetPublishInterval(expected)

	actual := configManager.PublishInterval()

	if expected != actual {
		t.Fail()
	}
}

func newConfigManager() LocalFileConfigManager {
	return LocalFileConfigManager{
		FilePath:               "./test_config.json",
		DefaultAppID:           "595068606",
		DefaultPublishInterval: 24 * time.Hour,
	}
}
