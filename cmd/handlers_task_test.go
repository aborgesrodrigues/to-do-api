package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	mock_svc "github.com/aborgesrodrigues/to-do-api/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	task := &common.Task{
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.addTask)

	errAddTask := errors.New("error inserting task")
	tests := map[string]struct {
		task         *common.Task
		svcError     error
		expectedResp string
	}{
		"success": {
			task:         task,
			svcError:     nil,
			expectedResp: `{"message":"Task Added"}`,
		},
		"fail": {
			task:         task,
			svcError:     errAddTask,
			expectedResp: `"error inserting task"`,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(task)
			assert.NoError(t, err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("POST", "/tasks", ioutil.NopCloser(&buf))

			svcInterface.
				EXPECT().
				AddTask(test.task).
				Return(test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				assert.Equal(t, http.StatusCreated, rr.Code)
			} else {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
			}
			assert.Equal(t, test.expectedResp, strings.TrimSpace(rr.Body.String()))
		})

	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	idTask := "0001"
	task := &common.Task{
		Id:          idTask,
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.updateTask)

	ctx := context.WithValue(context.Background(), taskIdCtx, idTask)

	errUpdateTask := errors.New("error updating task")
	tests := map[string]struct {
		task         *common.Task
		svcError     error
		expectedResp string
	}{
		"success": {
			task:         task,
			svcError:     nil,
			expectedResp: `{"message":"Task Updated"}`,
		},
		"fail": {
			task:         task,
			svcError:     errUpdateTask,
			expectedResp: `"error updating task"`,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(task)
			assert.NoError(t, err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("PUT", "/tasks/"+idTask, ioutil.NopCloser(&buf)).WithContext(ctx)

			svcInterface.
				EXPECT().
				UpdateTask(test.task).
				Return(test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				assert.Equal(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
			}
			assert.Equal(t, test.expectedResp, strings.TrimSpace(rr.Body.String()))
		})

	}
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	idTask := "0001"
	task := &common.Task{
		Id:          idTask,
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.getTask)

	ctx := context.WithValue(context.Background(), taskIdCtx, idTask)

	errGetTask := errors.New("error retrieving task")
	tests := map[string]struct {
		svcError     error
		expectedResp string
	}{
		// "success": {
		// 	svcError:     nil,
		// 	expectedResp: `{"id":"0001","user_id":"00001","description":"description 1","state":"to_do"}`,
		// },
		"fail": {
			svcError:     errGetTask,
			expectedResp: `"error retrieving task"`,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(task)
			assert.NoError(t, err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", "/tasks/"+idTask, ioutil.NopCloser(&buf)).WithContext(ctx)

			svcInterface.
				EXPECT().
				GetTask(idTask).
				Return(task, test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				assert.Equal(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
			}
			assert.Equal(t, test.expectedResp, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	idTask := "0001"

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.deleteTask)

	ctx := context.WithValue(context.Background(), taskIdCtx, idTask)

	errDeleteTask := errors.New("error deleting task")
	tests := map[string]struct {
		svcError     error
		expectedResp string
	}{
		"success": {
			svcError:     nil,
			expectedResp: `{"message":"Task Deleted"}`,
		},
		"fail": {
			svcError:     errDeleteTask,
			expectedResp: `"error deleting task"`,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("DELETE", "/tasks/"+idTask, nil).WithContext(ctx)

			svcInterface.
				EXPECT().
				DeleteTask(idTask).
				Return(test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				assert.Equal(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
			}
			assert.Equal(t, test.expectedResp, strings.TrimSpace(rr.Body.String()))
		})

	}
}

func TestListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	tasks := []common.Task{
		{
			Id:          "0001",
			UserId:      "00001",
			Description: "description 1",
			State:       "to_do",
		},
		{
			Id:          "0002",
			UserId:      "00002",
			Description: "description 2",
			State:       "to_do",
		},
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.listTasks)

	errGetTasks := errors.New("error retrieving tasks")
	tests := map[string]struct {
		tasks        []common.Task
		svcError     error
		expectedResp string
	}{
		"success": {
			tasks:        tasks,
			svcError:     nil,
			expectedResp: `[{"id":"0001","user_id":"00001","description":"description 1","state":"to_do"},{"id":"0002","user_id":"00002","description":"description 2","state":"to_do"}]`,
		},
		"fail": {
			svcError:     errGetTasks,
			expectedResp: `"error retrieving tasks"`,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", "/tasks/", nil)

			svcInterface.
				EXPECT().
				ListTasks().
				Return(tasks, test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				assert.Equal(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
			}
			assert.Equal(t, test.expectedResp, strings.TrimSpace(rr.Body.String()))
		})
	}
}
