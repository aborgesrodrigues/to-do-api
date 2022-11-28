package service

import (
	"errors"
	"testing"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	mock_db "github.com/aborgesrodrigues/to-do-api/internal/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errAddUser := errors.New("error inserting user")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}

	tests := map[string]struct {
		user         *common.User
		dbError      error
		expectedResp error
	}{
		"success": {
			user:         user,
			dbError:      nil,
			expectedResp: nil,
		},
		"fail": {
			user:         user,
			dbError:      errAddUser,
			expectedResp: errAddUser,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				AddUser(test.user).
				Return(test.dbError)

			err := svc.AddUser(test.user)
			assert.Equal(t, err, test.expectedResp)
		})

	}
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errAddUser := errors.New("error inserting user")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}

	tests := map[string]struct {
		user         *common.User
		dbError      error
		expectedResp error
	}{
		"success": {
			user:         user,
			dbError:      nil,
			expectedResp: nil,
		},
		"fail": {
			user:         user,
			dbError:      errAddUser,
			expectedResp: errAddUser,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				UpdateUser(test.user).
				Return(test.dbError)

			err := svc.UpdateUser(test.user)
			assert.Equal(t, err, test.expectedResp)
		})

	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errGetUser := errors.New("any error")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}

	tests := map[string]struct {
		id           string
		dbError      error
		dbUser       *common.User
		expectedResp *common.User
		expectedErr  error
	}{
		"success": {
			dbError:      nil,
			dbUser:       user,
			expectedResp: user,
			expectedErr:  nil,
		},
		"fail": {
			dbError:      errGetUser,
			dbUser:       nil,
			expectedResp: nil,
			expectedErr:  errGetUser,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				GetUser(test.id).
				Return(test.dbUser, test.dbError)

			user, err := svc.GetUser(test.id)
			assert.Equal(t, user, test.expectedResp)
			assert.Equal(t, err, test.expectedErr)
		})

	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errAddUser := errors.New("error inserting user")

	tests := map[string]struct {
		id           string
		dbError1     error
		dbError2     error
		expectedResp error
	}{
		"success": {
			id:           "0001",
			dbError1:     nil,
			dbError2:     nil,
			expectedResp: nil,
		},
		"fail1": {
			id:           "0001",
			dbError1:     errAddUser,
			dbError2:     nil,
			expectedResp: errAddUser,
		},
		"fail2": {
			id:           "0001",
			dbError1:     nil,
			dbError2:     errAddUser,
			expectedResp: errAddUser,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				DeleteUserTasks(test.id).
				Return(test.dbError1)

			if test.dbError1 == nil {
				dbInterface.
					EXPECT().
					DeleteUser(test.id).
					Return(test.dbError2)
			}

			err := svc.DeleteUser(test.id)
			assert.Equal(t, err, test.expectedResp)
		})
	}
}

func TestListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errGetUsers := errors.New("any error")
	users := []common.User{
		{
			Username: "username1",
			Name:     "User Name 1",
		},
		{
			Username: "username2",
			Name:     "User Name 2",
		},
	}

	tests := map[string]struct {
		dbError      error
		dbUsers      []common.User
		expectedResp []common.User
		expectedErr  error
	}{
		"success": {
			dbError:      nil,
			dbUsers:      users,
			expectedResp: users,
			expectedErr:  nil,
		},
		"fail": {
			dbError:      errGetUsers,
			dbUsers:      nil,
			expectedResp: nil,
			expectedErr:  errGetUsers,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				ListUsers().
				Return(test.dbUsers, test.dbError)

			user, err := svc.ListUsers()
			assert.Equal(t, user, test.expectedResp)
			assert.Equal(t, err, test.expectedErr)
		})

	}
}

func TestListUserTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errGetUsers := errors.New("any error")
	tasks := []common.Task{
		{
			UserId:      "111",
			Description: "description 1",
			State:       "to_do",
		},
		{
			UserId:      "222",
			Description: "description 2",
			State:       "to_do",
		},
	}

	tests := map[string]struct {
		id           string
		dbError      error
		dbTasks      []common.Task
		expectedResp []common.Task
		expectedErr  error
	}{
		"success": {
			dbError:      nil,
			dbTasks:      tasks,
			expectedResp: tasks,
			expectedErr:  nil,
		},
		"fail": {
			dbError:      errGetUsers,
			dbTasks:      nil,
			expectedResp: nil,
			expectedErr:  errGetUsers,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				ListUserTasks(test.id).
				Return(test.dbTasks, test.dbError)

			users, err := svc.ListUserTasks(test.id)
			assert.Equal(t, users, test.expectedResp)
			assert.Equal(t, err, test.expectedErr)
		})

	}
}
