package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
)

func (hdl *handlerTestSuite) TestAddUser() {
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.AddUser)

	errAddUser := errors.New("error inserting user")
	tests := map[string]struct {
		user         *common.User
		svcError     error
		expectedResp string
	}{
		"success": {
			user:         user,
			svcError:     nil,
			expectedResp: `{"id":"","username":"username1","name":"User Name 1"}`,
		},
		"fail": {
			user:         user,
			svcError:     errAddUser,
			expectedResp: `"error inserting user"`,
		},
	}

	for index, test := range tests {
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(user)
			hdl.Assert().NoError(err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("POST", "/users", io.NopCloser(&buf))
			id := uuid.New().String()
			responseUser := &common.User{
				Id:       id,
				Username: test.user.Username,
				Name:     test.user.Name,
			}

			// set up service mock
			hdl.getService().
				AddUser(test.user).
				Return(responseUser, test.svcError)

			handler.ServeHTTP(rr, req)

			if test.svcError == nil {
				hdl.Assert().Equal(http.StatusCreated, rr.Code)
				err = json.NewDecoder(rr.Body).Decode(responseUser)
				hdl.Assert().NoError(err)
				hdl.Assert().Equal(id, responseUser.Id)
				hdl.Assert().NotEmpty(responseUser.Id)
			} else {
				hdl.Assert().Equal(http.StatusInternalServerError, rr.Code)
				hdl.Assert().Equal(test.expectedResp, strings.TrimSpace(rr.Body.String()))
			}

		})

	}
}

func (hdl *handlerTestSuite) TestUpdateUser() {
	idUser := "0001"
	user := &common.User{
		Id:       idUser,
		Username: "username1",
		Name:     "User Name 1",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.UpdateUser)

	ctx := context.WithValue(context.Background(), userIdCtx, idUser)

	errUpdateUser := errors.New("error updating user")
	tests := map[string]struct {
		user         *common.User
		svcError     error
		expectedResp string
	}{
		"success": {
			user:         user,
			svcError:     nil,
			expectedResp: `{"id":"","username":"username1","name":"User Name 1"}`,
		},
		"fail": {
			user:         user,
			svcError:     errUpdateUser,
			expectedResp: `"error updating user"`,
		},
	}

	for index, test := range tests {
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(user)
			hdl.Assert().NoError(err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("PUT", "/users/"+idUser, io.NopCloser(&buf)).WithContext(ctx)

			// set up service mock
			hdl.getService().
				UpdateUser(test.user).
				Return(test.user, test.svcError)

			handler.ServeHTTP(rr, req)

			if test.svcError == nil {
				hdl.Assert().Equal(http.StatusOK, rr.Code)
				err = json.NewDecoder(rr.Body).Decode(&test.user)
				hdl.Assert().NoError(err)
			} else {
				hdl.Assert().Equal(http.StatusInternalServerError, rr.Code)
				hdl.Assert().Equal(test.expectedResp, strings.TrimSpace(rr.Body.String()))
			}

		})

	}
}

func (hdl *handlerTestSuite) TestGetUser() {
	idUser := "0001"
	user := &common.User{
		Id:       idUser,
		Username: "username1",
		Name:     "User Name 1",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.GetUser)

	ctx := context.WithValue(context.Background(), userIdCtx, idUser)

	errGetUser := errors.New("error retrieving user")
	tests := map[string]struct {
		svcError     error
		expectedResp string
	}{
		"success": {
			svcError:     nil,
			expectedResp: `{"id":"0001","username":"username1","name":"User Name 1"}`,
		},
		"fail": {
			svcError:     errGetUser,
			expectedResp: `"error retrieving user"`,
		},
	}

	for index, test := range tests {
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(user)
			hdl.Assert().NoError(err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", "/users/"+idUser, io.NopCloser(&buf)).WithContext(ctx)

			// set up service mock
			hdl.getService().
				GetUser(idUser).
				Return(user, test.svcError)

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

func (hdl *handlerTestSuite) TestDeleteUser() {
	idUser := "0001"

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.DeleteUser)

	ctx := context.WithValue(context.Background(), userIdCtx, idUser)

	errDeleteUser := errors.New("error deleting user")
	tests := map[string]struct {
		svcError     error
		expectedResp string
	}{
		"success": {
			svcError:     nil,
			expectedResp: `{"message":"User Deleted"}`,
		},
		"fail": {
			svcError:     errDeleteUser,
			expectedResp: `"error deleting user"`,
		},
	}

	for index, test := range tests {
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("DELETE", "/users/"+idUser, nil).WithContext(ctx)

			// set up service mock
			hdl.getService().
				DeleteUser(idUser).
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

func (hdl *handlerTestSuite) TestListUsers() {
	users := []common.User{
		{
			Id:       "0001",
			Username: "username1",
			Name:     "User Name 1",
		},
		{
			Id:       "0002",
			Username: "username2",
			Name:     "User Name 2",
		},
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.ListUsers)

	errGetUsers := errors.New("error retrieving users")
	tests := map[string]struct {
		users        []common.User
		svcError     error
		expectedResp string
	}{
		"success": {
			users:        users,
			svcError:     nil,
			expectedResp: `[{"id":"0001","username":"username1","name":"User Name 1"},{"id":"0002","username":"username2","name":"User Name 2"}]`,
		},
		"fail": {
			svcError:     errGetUsers,
			expectedResp: `"error retrieving users"`,
		},
	}

	for index, test := range tests {
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", "/users/", nil)

			hdl.getService().
				ListUsers().
				Return(users, test.svcError)

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

func (hdl *handlerTestSuite) TestListUserTasks() {
	idUser := "0001"
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
	ctx := context.WithValue(context.Background(), userIdCtx, idUser)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.handler.ListUserTasks)

	errGetUsers := errors.New("error retrieving users")
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
			svcError:     errGetUsers,
			expectedResp: `"error retrieving users"`,
		},
	}

	for index, test := range tests {
		hdl.Run(index, func() {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s/tasks", idUser), nil).WithContext(ctx)

			hdl.getService().
				ListUserTasks(idUser).
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
