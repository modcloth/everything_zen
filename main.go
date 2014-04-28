package main

import (
	"fmt"
	"log"
	"os"

	"github.com/modcloth/everything_zen/tracker"
)

func main() {

	agileToken := os.Getenv("AGILE_TOKEN")
	agileProject := os.Getenv("AGILE_PROJECTID")
	trackerToken := os.Getenv("TRACKER_TOKEN")
	trackerProject := os.Getenv("TRACKER_PROJECTID")

	agileStories := GetAgileStories(agileToken, agileProject)
	filtered := FilterNonTracker(agileStories)
	tasks := AgileToTask(filtered)
	AddOrUpdateTasks(tasks, trackerToken, trackerProject)
}

func AddOrUpdateTasks(tasks <-chan TrackerTask, token, projectID string) {
	tracker, err := tracker.NewTracker(token, projectID)
	LogAndQuit(err)

	for task := range tasks {
		fmt.Printf("assigning agile task %s to pivotal story %d\n", task.AgileURL, task.StoryID)
		tracker.SaveTask(task.StoryID, task.AgileURL, task.IsComplete)
	}
}

func LogAndQuit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
