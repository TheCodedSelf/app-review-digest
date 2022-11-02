package main

import (
	"fmt"
	"time"
)

func main() {
	var configManager ConfigManager = NewLocalFileConfigManager()

	publishTime := configManager.GetPublishTime()

	publisher := NewPublisher(configManager)
	f := func() {
		currentTime := time.Now()
		fmt.Printf("Publishing digest if necessary at %v", currentTime)
		publisher.PublishLatest(currentTime)
	}
	ScheduleJob(publishTime, f)
}
