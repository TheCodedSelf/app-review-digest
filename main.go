/*
main.go

See readme.md for guidance
*/

package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	// Command line flags
	hourFlag := flag.Int("hour", -1, "Optional. Change the hour to schedule digests between 0 and 23.")
	minuteFlag := flag.Int("minute", -1, "Optional. Change the minute to schedule digests between 0 and 59.")
	intervalFlag := flag.Int("interval", -1, "Optional. Change the interval (how many days of reviews) of the digest publication. Specified in days. Default 1.")
	appIDFlag := flag.String("app", "", "Optional. Configure the Apple app ID to poll.")
	nowFlag := flag.Bool("now", false, "Optional. Try publish a new digest and then exit. Returns the last digest if published in the last 24 hours")
	flag.Parse()

	var configManager ConfigManager = NewLocalFileConfigManager()
	publishTime := configManager.PublishTime()

	// If configured, set the scheduler's run time
	if *hourFlag >= 0 || *minuteFlag >= 0 {
		if *hourFlag >= 0 {
			publishTime.Hour = *hourFlag
		}
		if *minuteFlag >= 0 {
			publishTime.Minute = *minuteFlag
		}
		fmt.Printf("Setting new publish time to %v\n", publishTime)
		configManager.SetPublishTime(publishTime)
	}

	// If configured, set the length (interval) of the digest
	if *intervalFlag > 0 {
		interval := time.Duration(*intervalFlag) * time.Hour * 24
		configManager.SetPublishInterval(interval)
	}

	// If configured, set the app ID of an app on the Apple App Store
	if *appIDFlag != "" {
		fmt.Printf("Setting new app ID to %s\n", *appIDFlag)
		configManager.SetAppID(*appIDFlag)
	}

	publisher := NewPublisher(configManager)

	// If configured, publish the latest digest and then exit
	if *nowFlag {
		path := publisher.PublishLatest(time.Now())
		fmt.Printf("Latest digest: %s", path)
		return
	}

	f := func() {
		currentTime := time.Now()
		fmt.Printf("Publishing digest if necessary at %v", currentTime)
		publisher.PublishLatest(currentTime)
	}

	// Schedule a job to publish the digest at the configured time
	ScheduleJob(publishTime, f)
}
