package main

import (
	"github.com/modcloth/everything_zen/agile_zen"
)

type TrackerTask struct {
	StoryID    int
	AgileURL   string
	IsComplete bool
}

type AgileStory struct {
	URL        string
	TrackerIDs []int
	IsComplete bool
}

func AgileToTask(stories <-chan AgileStory) <-chan TrackerTask {
	c := make(chan TrackerTask)

	go func() {
		defer close(c)
		for story := range stories {
			for _, trackerID := range story.TrackerIDs {
				c <- TrackerTask{
					StoryID:    trackerID,
					AgileURL:   story.URL,
					IsComplete: story.IsComplete,
				}
			}
		}

	}()
	return c
}

func FilterNonTracker(stories <-chan AgileStory) <-chan AgileStory {
	c := make(chan AgileStory)
	go func() {
		defer close(c)
		for story := range stories {
			if len(story.TrackerIDs) > 0 {
				c <- story
			}
		}
	}()
	return c
}

func GetAgileStories(token, projectID string) <-chan AgileStory {
	c := make(chan AgileStory)
	agile, err := agile_zen.NewAgileZen(token, projectID)
	LogAndQuit(err)

	stories, err := agile.Stories()
	LogAndQuit(err)

	go func() {
		defer close(c)
		for _, story := range stories {
			c <- AgileStory{URL: story.UrlForStory(),
				TrackerIDs: story.TrackerStories(),
				IsComplete: story.IsComplete(),
			}
		}
	}()
	return c
}
