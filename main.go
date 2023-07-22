package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/csmith/envflag"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

var (
	apiToken = flag.String("api-token", "", "Personal API token for Todoist")
)

func main() {
	envflag.Parse()

	tasks := getTasks()

	for i := range tasks {
		if tasks[i].Completed || tasks[i].Due.Recurring {
			continue
		}

		label := ""
		age := time.Now().Sub(tasks[i].Created)
		if age < 14*24*time.Hour {
			continue
		} else if age < 56*24*time.Hour {
			label = "age-weeks"
		} else if age < 365*24*time.Hour {
			label = "age-months"
		} else {
			label = "age-years"
		}

		if !slices.Contains(tasks[i].Labels, label) {
			newLabels := []string{label}
			for _, l := range tasks[i].Labels {
				if !strings.HasPrefix(l, "age-") {
					newLabels = append(newLabels, l)
				}
			}
			updateTask(tasks[i].ID, newLabels)
		}
	}
}

type Task struct {
	ID        string    `json:"id"`
	Labels    []string  `json:"labels"`
	Created   time.Time `json:"created_at"`
	Completed bool      `json:"is_completed"`
	Due       struct {
		Recurring bool `json:"is_recurring"`
	} `json:"due"`
}

func getTasks() []Task {
	req, err := http.NewRequest("GET", "https://api.todoist.com/rest/v2/tasks", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var tasks []Task

	err = json.Unmarshal(b, &tasks)
	if err != nil {
		panic(err)
	}

	return tasks
}

func updateTask(id string, labels []string) {
	body := struct {
		Labels []string `json:"labels"`
	}{labels}

	j, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.todoist.com/rest/v2/tasks/%s", id), bytes.NewReader(j))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *apiToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Request-Id", uuid.NewString())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected status code: %d", res.StatusCode))
	}
}
