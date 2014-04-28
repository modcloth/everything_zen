package agile_zen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	BaseURL = "https://agilezen.com/api/v1/projects"
)

type AgileZen struct {
	apiToken, url, projectID string
}

type Page struct {
	Page       int     `json:"page"`
	PageSize   int     `json:"pageSize"`
	TotalPages int     `json:"totalPages"`
	TotalItems int     `json:"totalItems"`
	Items      []Story `json:"items"`
}

type Story struct {
	ID       int     `json"id"`
	Text     string  `json"text"`
	Size     string  `json"size"`
	Color    string  `json"color"`
	Priority string  `json"priority"`
	Status   string  `json"ready"`
	Project  Project `json"project"`
	Phase    Phase   `json"phase"`
	Creator  User    `json"creator"`
	Owner    User    `json"owner"`
	Tags     []Tag   `json"tags"`
	URL      string  `json:"-"`
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Phase struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewAgileZen(token, projectID string) (*AgileZen, error) {
	if token == "" || projectID == "" {
		return nil, errors.New("AgileZen requires that both Token and projectID be specified")
	}

	url := fmt.Sprintf("%s/%s", BaseURL, projectID)
	return &AgileZen{token, url, projectID}, nil
}
func (az *AgileZen) Stories() ([]Story, error) {
	items := make([]Story, 0)
	currentPage := 1

	for {
		page := Page{}
		url := fmt.Sprintf("%s/stories?page=%d&with=tags", az.url, currentPage)
		if err := az.doGet(url, &page); err != nil {
			return nil, err
		}
		items = append(items, page.Items...)
		if currentPage == page.TotalPages {
			break
		}
		currentPage += 1
	}
	return items, nil
}

func (s *Story) UrlForStory() string {
	return fmt.Sprintf("https://agilezen.com/project/%d/story/%d", s.Project.ID, s.ID)
}

func (s *Story) TrackerStories() []int {
	ids := make([]int, 0)
	for _, tag := range s.Tags {
		if id, err := strconv.Atoi(tag.Name); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func (s *Story) IsComplete() bool {
	return s.Phase.Name == "Complete" || s.Phase.Name == "Archive"
}

func (t *AgileZen) doGet(url string, item interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-Zen-ApiKey", t.apiToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	res.Body.Close()
	err = json.Unmarshal(buffer, &item)
	if err != nil {
		return err
	}
	return nil
}
