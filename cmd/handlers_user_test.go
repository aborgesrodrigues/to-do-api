package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.addUser)

	errAddUser := errors.New("error inserting user")
	tests := map[string]struct {
		user         *common.User
		svcError     error
		expectedResp string
	}{
		"success": {
			user:         user,
			svcError:     nil,
			expectedResp: `{"message":"User Added"}`,
		},
		"fail": {
			user:         user,
			svcError:     errAddUser,
			expectedResp: `"error inserting user"`,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(user)
			assert.NoError(t, err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("POST", "/users", ioutil.NopCloser(&buf))

			svcInterface.
				EXPECT().
				AddUser(test.user).
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

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	idUser := "0001"
	user := &common.User{
		Id:       idUser,
		Username: "username1",
		Name:     "User Name 1",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.updateUser)

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
			expectedResp: `{"message":"User Updated"}`,
		},
		"fail": {
			user:         user,
			svcError:     errUpdateUser,
			expectedResp: `"error updating user"`,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(user)
			assert.NoError(t, err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("PUT", "/users/"+idUser, ioutil.NopCloser(&buf)).WithContext(ctx)

			svcInterface.
				EXPECT().
				UpdateUser(test.user).
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

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	idUser := "0001"
	user := &common.User{
		Id:       idUser,
		Username: "username1",
		Name:     "User Name 1",
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.getUser)

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
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(user)
			assert.NoError(t, err)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", "/users/"+idUser, ioutil.NopCloser(&buf)).WithContext(ctx)

			svcInterface.
				EXPECT().
				GetUser(idUser).
				Return(user, test.svcError)

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

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

	idUser := "0001"

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	handler := http.HandlerFunc(hdl.deleteUser)

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
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("DELETE", "/users/"+idUser, nil).WithContext(ctx)

			svcInterface.
				EXPECT().
				DeleteUser(idUser).
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

func TestListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

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
	handler := http.HandlerFunc(hdl.listUsers)

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
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", "/users/", nil)

			svcInterface.
				EXPECT().
				ListUsers().
				Return(users, test.svcError)

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

func TestListUserTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creates service mock
	svcInterface := mock_svc.NewMockSVCInterface(ctrl)
	hdl.svc = svcInterface

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
	handler := http.HandlerFunc(hdl.listUserTasks)

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
		t.Run(index, func(t *testing.T) {
			rr := httptest.NewRecorder()

			// Create a request to pass to our handler.
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s/tasks", idUser), nil).WithContext(ctx)

			svcInterface.
				EXPECT().
				ListUserTasks(idUser).
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
