package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	BaseURL = "https://www.pivotaltracker.com/services/v5/projects"
)

type Tracker struct {
	apiToken, projectId string
	url                 string
}

func NewTracker(apiToken, projectID string) (*Tracker, error) {
	if apiToken == "" || projectID == "" {
		return nil, errors.New("Tracker requires that both Token and projectID be specified")
	}
	return &Tracker{apiToken: apiToken, url: fmt.Sprintf("%s/%s/", BaseURL, projectID)}, nil
}

// https://www.pivotaltracker.com/help/api/rest/v5#story_resource
type Story struct {
	Estimate               float64   `json:"estimate"`
	OwnerIDs               []int     `json:"owner_ids"`
	Labels                 []Label   `json:"labels"`
	FollowerIDs            []int     `json:"follower_ids"`
	CommentIDs             []int     `json:"comment_ids"`
	CreatedAt              time.Time `json:"created_at"`
	ExternalID             string    `json:"external_id"`
	Description            string    `json:"description"`
	UpdatedAt              time.Time `json:"updated_at"`
	StoryType              string    `json:"string"`
	RequestedByID          int       `json:"requested_by_id"`
	URL                    string    `json:"url"`
	ProjectID              int       `json:"project_id"`
	AcceptedAt             time.Time `json:"accepted_at"`
	IntegrationID          int       `json:"integration_id"`
	CurrentState           string    `json:"current_state"`
	PlannedIterationNumber int       `json:"planned_iteration_number"`
	ID                     int       `json:"id"`
	Name                   string    `json:"name"`
	Tasks                  []Task    `json:"tasks"`
}

func (t *Tracker) AddTask(storyID int, text string, isComplete bool) error {

	url := fmt.Sprintf("%s/stories/%d/tasks", t.url, storyID)
	body := strings.NewReader(fmt.Sprintf(`{"description":"%s", "complete": %t}`, text, isComplete))

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("X-TrackerToken", t.apiToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	return nil
}

func (t *Tracker) UpdateTask(storyID, taskID int, isComplete bool) error {

	url := fmt.Sprintf("%s/stories/%d/tasks/%d", t.url, storyID, taskID)
	body := strings.NewReader(fmt.Sprintf(`{"complete": %t}`, isComplete))

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("X-TrackerToken", t.apiToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	return nil
}

func (t *Tracker) SaveTask(storyID int, text string, isComplete bool) {
	tasks := t.GetTasks(storyID)

	for _, task := range tasks {
		if task.Description == text {
			if task.Complete != isComplete {
				t.UpdateTask(storyID, task.ID, isComplete)
			}
			return
		}
	}
	t.AddTask(storyID, text, isComplete)
}

func (t *Tracker) GetTasks(storyID int) []Task {
	url := fmt.Sprintf("%s/stories/%d/tasks", t.url, storyID)
	var tasks []Task
	t.doGet(url, &tasks)
	return tasks
}

type Label struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	Complete    bool      `json:"complete"`
	Position    int       `json:"position"`
	StoryID     int       `json:"story_id"`
}

func (t *Tracker) doGet(url string, item interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-TrackerToken", t.apiToken)

	client := &http.Client{}
	res, err := client.Do(req)
	Kaboom(err)

	buffer, err := ioutil.ReadAll(res.Body)
	Kaboom(err)

	res.Body.Close()
	err = json.Unmarshal(buffer, &item)
	if err != nil {
		return err
	}
	return nil
}

func Kaboom(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
