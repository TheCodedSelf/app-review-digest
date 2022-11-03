/*
job_scheduler.go

Allows scheduling and running a job at a particular time interval.
*/
package main

import (
	"fmt"
	"log"
	"time"
)

// For demonstration purposes, consider changing the interval period to once per minute:
// const JOB_INTERVAL_PERIOD time.Duration = 1 * time.Minute

// Or, leave at 24 hours.
const JOB_INTERVAL_PERIOD time.Duration = 24 * time.Hour

type JobScheduler struct {
	timer *time.Timer
	at    TimeOfDay
	job   func()
}

// Creates and runs a scheduler
// Takes a time of day for the next execution and a closure to execute
func ScheduleJob(at TimeOfDay, job func()) {
	scheduler := JobScheduler{}
	scheduler.timer = time.NewTimer(nextTickDuration(at))
	scheduler.at = at
	scheduler.job = job
	scheduler.run()
}

// Runs the scheduler's closure and updates the job ticker
func (j JobScheduler) run() {
	for {
		<-j.timer.C
		fmt.Println(time.Now(), "- job scheduler executing")
		j.job()
		j.updateJobTicker()
	}
}

// Resets the scheduler's timer to the next time the scheduler should run
func (j JobScheduler) updateJobTicker() {
	nextTick := nextTickDuration(j.at)
	fmt.Printf("Updating job ticker to run at %d", nextTick)
	j.timer.Reset(nextTick)
}

// Calculates the next time at which the job should run
func nextTickDuration(at TimeOfDay) time.Duration {
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), at.Hour, at.Minute, 0, 0, time.Local)
	for nextTick.Before(now) {
		nextTick = nextTick.Add(JOB_INTERVAL_PERIOD)
	}
	duration := nextTick.Sub(time.Now())
	fmt.Printf("Next job will execute at %v (%v from now)\n", nextTick, duration)
	if duration <= 0 {
		log.Fatal("Scheduled job at bad duration.")
	}
	return duration
}
