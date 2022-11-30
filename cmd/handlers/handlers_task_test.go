package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
)

func (hdl *handlerTestSuite) TestAddTask() {
	task := &common.Task{
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.AddTask)

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
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(task)
			hdl.Assert().NoError(err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("POST", "/tasks", ioutil.NopCloser(&buf))
			id := uuid.New().String()
			responseTask := &common.Task{
				Id:          id,
				Description: test.task.Description,
				UserId:      test.task.UserId,
				State:       test.task.State,
			}

			hdl.getService().
				AddTask(test.task).
				Return(responseTask, test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				hdl.Assert().Equal(http.StatusCreated, rr.Code)
				err = json.NewDecoder(rr.Body).Decode(responseTask)
				hdl.Assert().NoError(err)
				hdl.Assert().Equal(id, responseTask.Id)
				hdl.Assert().NotEmpty(responseTask.Id)
			} else {
				hdl.Assert().Equal(http.StatusInternalServerError, rr.Code)
				hdl.Assert().Equal(test.expectedResp, strings.TrimSpace(rr.Body.String()))
			}
		})

	}
}

func (hdl *handlerTestSuite) TestUpdateTask() {
	idTask := "0001"
	task := &common.Task{
		Id:          idTask,
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.UpdateTask)

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
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(task)
			hdl.Assert().NoError(err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("PUT", "/tasks/"+idTask, ioutil.NopCloser(&buf)).WithContext(ctx)

			hdl.getService().
				UpdateTask(test.task).
				Return(test.task, test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				hdl.Assert().Equal(http.StatusOK, rr.Code)
				err = json.NewDecoder(rr.Body).Decode(&test.task)
				hdl.Assert().NoError(err)
			} else {
				hdl.Assert().Equal(http.StatusInternalServerError, rr.Code)
				hdl.Assert().Equal(test.expectedResp, strings.TrimSpace(rr.Body.String()))
			}
		})

	}
}

func (hdl *handlerTestSuite) TestGetTask() {
	idTask := "0001"
	task := &common.Task{
		Id:          idTask,
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.GetTask)

	ctx := context.WithValue(context.Background(), taskIdCtx, idTask)

	errGetTask := errors.New("error retrieving task")
	tests := map[string]struct {
		svcError     error
		expectedResp string
	}{
		"success": {
			svcError:     nil,
			expectedResp: `{"id":"0001","user_id":"00001","description":"description 1","state":"to_do"}`,
		},
		"fail": {
			svcError:     errGetTask,
			expectedResp: `"error retrieving task"`,
		},
	}

	for index, test := range tests {
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(task)
			hdl.Assert().NoError(err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", "/tasks/"+idTask, ioutil.NopCloser(&buf)).WithContext(ctx)

			hdl.getService().
				GetTask(idTask).
				Return(task, test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				hdl.Assert().Equal(http.StatusOK, rr.Code)
			} else {
				hdl.Assert().Equal(http.StatusInternalServerError, rr.Code)
			}
			hdl.Assert().Equal(test.expectedResp, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func (hdl *handlerTestSuite) TestDeleteTask() {
	idTask := "0001"

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.DeleteTask)

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
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("DELETE", "/tasks/"+idTask, nil).WithContext(ctx)

			hdl.getService().
				DeleteTask(idTask).
				Return(test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				hdl.Assert().Equal(http.StatusOK, rr.Code)
			} else {
				hdl.Assert().Equal(http.StatusInternalServerError, rr.Code)
			}
			hdl.Assert().Equal(test.expectedResp, strings.TrimSpace(rr.Body.String()))
		})

	}
}

func (hdl *handlerTestSuite) TestListTasks() {
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
	handler := http.HandlerFunc(hdl.handler.ListTasks)

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
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", "/tasks/", nil)

			hdl.getService().
				ListTasks().
				Return(tasks, test.svcError)

			handler.ServeHTTP(rr, req)
			if test.svcError == nil {
				hdl.Assert().Equal(http.StatusOK, rr.Code)
			} else {
				hdl.Assert().Equal(http.StatusInternalServerError, rr.Code)
			}
			hdl.Assert().Equal(test.expectedResp, strings.TrimSpace(rr.Body.String()))
		})
	}
}
