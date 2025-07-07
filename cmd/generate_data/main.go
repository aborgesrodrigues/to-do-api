package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
)

func main() {
	inputChan := make(chan common.User, 100)
	var wg sync.WaitGroup

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(inputChan)
		}()
	}

	go func() {
		for i := range 10000 {
			inputChan <- common.User{
				Username: fmt.Sprintf("user%d", i+1),
				Name:     fmt.Sprintf("User %d", i+1),
			}
		}

		close(inputChan)
	}()

	wg.Wait()
	log.Print("Finished")
}

func worker(input chan common.User) {
	adminUser := common.User{
		Username: "admin",
		Name:     "Admin",
	}
	response, err := doRequest[map[string]any](http.MethodPost, "http://localhost:8080/token", "", adminUser)
	if err != nil {
		log.Fatalf("Error requesting: %v", err)
	}

	bearerToken := (*response)["access_token"].(string)
	for user := range input {
		response, err := doRequest[map[string]any](http.MethodPost, "http://localhost:8080/users", bearerToken, user)
		if err != nil {
			log.Fatalf("Error requesting: %v", err)
		}

		user := (*response)["user"].(map[string]any)
		for i := range 20 {
			task := common.Task{
				UserId:      user["id"].(string),
				Description: fmt.Sprintf("Task %d", i+1),
				State:       common.TaskState("to_do"),
			}

			_, err := doRequest[map[string]any](http.MethodPost, "http://localhost:8080/tasks", bearerToken, task)
			if err != nil {
				log.Fatalf("Error requesting: %v", err)
			}
		}

		log.Printf("Added user %s, %s", user["id"], user["username"])
	}
}

func doRequest[T any](method, url, bearerToken string, input any) (*T, error) {
	bInput, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bInput))
	if err != nil {
		return nil, err
	}

	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	var output T
	err = json.Unmarshal(respBody, &output)
	if err != nil {
		log.Fatalf("Error unmarshalling: %#v %v", string(respBody), err)
	}

	return &output, nil
}
