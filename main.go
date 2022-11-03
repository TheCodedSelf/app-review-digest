package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	hourFlag := flag.Int("hour", -1, "Optional. Change the hour to schedule newsletter between 0 and 23.")
	minuteFlag := flag.Int("minute", -1, "Optional. Change the minute to schedule newsletter between 0 and 59.")
	appIDFlag := flag.String("app", "", "Optional. Configure the Apple app ID to poll.")
	nowFlag := flag.Bool("now", false, "Optional. Try publish a new newsletter and then exit.")
	flag.Parse()

	var configManager ConfigManager = NewLocalFileConfigManager()
	publishTime := configManager.PublishTime()

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

	if *appIDFlag != "" {
		fmt.Printf("Setting new app ID to %s\n", *appIDFlag)
		configManager.SetAppID(*appIDFlag)
	}

	publisher := NewPublisher(configManager)

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
	ScheduleJob(publishTime, f)
}
