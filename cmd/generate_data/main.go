package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
)

var inputChan chan (common.User)

func main() {
	logger, _ := zap.NewDevelopment()
	inputChan = make(chan common.User, 100)
	var wg sync.WaitGroup

	adminUser := common.User{
		Username: "admin",
		Name:     "Admin",
	}
	response, err := doRequest[map[string]any](http.MethodPost, "http://localhost:8080/token", "", adminUser, logger)
	if err != nil {
		logger.Fatal("Error requesting", zap.Error(err))
	}
	adminBearerToken := (*response)["access_token"].(string)

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(logger, adminBearerToken)
		}()
	}

	go func() {
		for i := range 10000 {
			inputChan <- common.User{
				Username: fmt.Sprintf("user%d", i+1),
				Name:     fmt.Sprintf("User %d", i+1),
				Password: fmt.Sprintf("password%d", i+1),
			}
		}

		close(inputChan)
	}()

	wg.Wait()
	logger.Info("Finished")
}

func worker(logger *zap.Logger, adminBearerToken string) {
	for userInput := range inputChan {
		response, err := doRequest[map[string]any](http.MethodPost, "http://localhost:8080/users", adminBearerToken, userInput, logger)
		if err != nil {
			logger.Fatal("Error requesting", zap.Error(err))
		}

		bearerToken := (*response)["access_token"].(string)

		user := (*response)["user"].(map[string]any)
		for i := range 20 {
			task := common.Task{
				UserId:      user["id"].(string),
				Description: fmt.Sprintf("Task %d", i+1),
				State:       common.TaskState("to_do"),
			}

			_, err := doRequest[map[string]any](http.MethodPost, fmt.Sprintf("http://localhost:8080/users/%s/tasks", user["id"].(string)), bearerToken, task, logger)
			if err != nil {
				logger.Fatal("Error requesting", zap.Error(err))
			}
		}

		logger.Info("Added user", zap.String("id", user["id"].(string)), zap.String("username", user["username"].(string)))
	}
}

func doRequest[T any](method, url, bearerToken string, input any, logger *zap.Logger) (*T, error) {
	bInput, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, io.NopCloser(bytes.NewReader(bInput)))
	if err != nil {
		return nil, err
	}

	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Error doing request", zap.String("method", method), zap.String("url", url), zap.String("input", string(bInput)), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	var output T
	err = json.Unmarshal(respBody, &output)
	if err != nil {
		logger.Fatal("Error unmarshalling", zap.String("respBody", string(respBody)), zap.Error(err))
	}

	return &output, nil
}
