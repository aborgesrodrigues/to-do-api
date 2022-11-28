package service

import (
	"errors"
	"testing"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	mock_db "github.com/aborgesrodrigues/to-do-api/internal/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errAddTask := errors.New("error inserting task")
	task := &common.Task{
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	tests := map[string]struct {
		task         *common.Task
		dbError      error
		expectedResp error
	}{
		"success": {
			task:         task,
			dbError:      nil,
			expectedResp: nil,
		},
		"fail": {
			task:         task,
			dbError:      errAddTask,
			expectedResp: errAddTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				AddTask(test.task).
				Return(test.dbError)

			err := svc.AddTask(test.task)
			assert.Equal(t, err, test.expectedResp)
		})

	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errAddTask := errors.New("error inserting task")
	task := &common.Task{
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	tests := map[string]struct {
		task         *common.Task
		dbError      error
		expectedResp error
	}{
		"success": {
			task:         task,
			dbError:      nil,
			expectedResp: nil,
		},
		"fail": {
			task:         task,
			dbError:      errAddTask,
			expectedResp: errAddTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				UpdateTask(test.task).
				Return(test.dbError)

			err := svc.UpdateTask(test.task)
			assert.Equal(t, err, test.expectedResp)
		})

	}
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errGetTask := errors.New("any error")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}
	task := &common.Task{
		UserId:      "00001",
		Description: "description 1",
		State:       "to_do",
	}

	tests := map[string]struct {
		id           string
		dbError1     error
		dbError2     error
		dbTask       *common.Task
		dbUser       *common.User
		expectedResp *common.Task
		expectedErr  error
	}{
		"success": {
			dbError1:     nil,
			dbError2:     nil,
			dbTask:       task,
			dbUser:       user,
			expectedResp: task,
			expectedErr:  nil,
		},
		"fail1": {
			dbError1:     errGetTask,
			dbError2:     nil,
			dbTask:       nil,
			dbUser:       nil,
			expectedResp: nil,
			expectedErr:  errGetTask,
		},
		"fail2": {
			dbError1:     nil,
			dbError2:     errGetTask,
			dbTask:       task,
			expectedResp: nil,
			expectedErr:  errGetTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				GetTask(test.id).
				Return(test.dbTask, test.dbError1)

			if test.dbError1 == nil {
				dbInterface.
					EXPECT().
					GetUser("00001").
					Return(&common.User{}, test.dbError2)
			}

			task, err := svc.GetTask(test.id)
			assert.Equal(t, task, test.expectedResp)
			assert.Equal(t, err, test.expectedErr)
		})

	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errAddTask := errors.New("error inserting task")

	tests := map[string]struct {
		id           string
		dbError      error
		expectedResp error
	}{
		"success": {
			id:           "0001",
			dbError:      nil,
			expectedResp: nil,
		},
		"fail": {
			id:           "0001",
			dbError:      errAddTask,
			expectedResp: errAddTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				DeleteTask(test.id).
				Return(test.dbError)

			err := svc.DeleteTask(test.id)
			assert.Equal(t, err, test.expectedResp)
		})
	}
}

func TestListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbInterface := mock_db.NewMockDBInterface(ctrl)
	svc, err := New(cfg)
	svc.db = dbInterface

	assert.NoError(t, err)

	errGetTasks := errors.New("any error")
	tasks := []common.Task{
		{
			UserId:      "00001",
			Description: "description 1",
			State:       "to_do",
		},
		{
			UserId:      "00002",
			Description: "description 2",
			State:       "to_do",
		},
	}

	tests := map[string]struct {
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
			dbError:      errGetTasks,
			dbTasks:      nil,
			expectedResp: nil,
			expectedErr:  errGetTasks,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			dbInterface.
				EXPECT().
				ListTasks().
				Return(test.dbTasks, test.dbError)

			tasks, err := svc.ListTasks()
			assert.Equal(t, tasks, test.expectedResp)
			assert.Equal(t, err, test.expectedErr)
		})

	}
}
