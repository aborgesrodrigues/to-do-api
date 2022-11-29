package db

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errAddTask := errors.New("error inserting task")
	task := &common.Task{
		UserId:      "0001",
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
			mockInsert := mock.ExpectExec("INSERT INTO public.task").WithArgs(test.task.Id, test.task.UserId, test.task.Description, test.task.State)
			if test.dbError == nil {
				mockInsert.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockInsert.WillReturnError(test.dbError)
			}

			err := db.AddTask(test.task)
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

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errAddTask := errors.New("error updating task")
	task := &common.Task{
		UserId:      "0001",
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
			mockUpdate := mock.ExpectExec("UPDATE public.task").WithArgs(test.task.UserId, test.task.Description, test.task.State, test.task.Id)
			if test.dbError == nil {
				mockUpdate.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockUpdate.WillReturnError(test.dbError)
			}

			err := db.UpdateTask(test.task)
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

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errGetTask := errors.New("any error")
	task := &common.Task{
		UserId:      "0001",
		Description: "description 1",
		State:       "to_do",
	}
	rowTask := sqlmock.NewRows([]string{"id", "user_id", "description", "state"}).AddRow("", task.UserId, task.Description, task.State)

	tests := map[string]struct {
		id           string
		dbError      error
		dbRowTask    *sqlmock.Rows
		expectedResp *common.Task
		expectedErr  error
	}{
		"success": {
			dbError:      nil,
			dbRowTask:    rowTask,
			expectedResp: task,
			expectedErr:  nil,
		},
		"fail": {
			dbError:      errGetTask,
			dbRowTask:    nil,
			expectedResp: nil,
			expectedErr:  errGetTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			mockGet := mock.ExpectQuery("SELECT id, user_id, description, state FROM public.task").WithArgs(test.id)
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowTask)
			} else {
				mockGet.WillReturnError(errGetTask)
			}

			task, err := db.GetTask(test.id)
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

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errAddTask := errors.New("error deleting task")

	tests := map[string]struct {
		id           string
		dbError      error
		expectedResp error
	}{
		"success": {
			dbError:      nil,
			expectedResp: nil,
		},
		"fail": {
			dbError:      errAddTask,
			expectedResp: errAddTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			mockUpdate := mock.ExpectExec("DELETE FROM public.task").WithArgs(test.id)
			if test.dbError == nil {
				mockUpdate.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockUpdate.WillReturnError(test.dbError)
			}

			err := db.DeleteTask(test.id)
			assert.Equal(t, err, test.expectedResp)
		})

	}
}

func TestDeleteUserTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errAddTask := errors.New("error deleting task")

	tests := map[string]struct {
		id           string
		dbError      error
		expectedResp error
	}{
		"success": {
			dbError:      nil,
			expectedResp: nil,
		},
		"fail": {
			dbError:      errAddTask,
			expectedResp: errAddTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			mockUpdate := mock.ExpectExec("DELETE FROM public.task").WithArgs(test.id)
			if test.dbError == nil {
				mockUpdate.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockUpdate.WillReturnError(test.dbError)
			}

			err := db.DeleteUserTasks(test.id)
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

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errGetTask := errors.New("any error")
	listTasks := []common.Task{
		{
			UserId:      "0001",
			Description: "description 1",
			State:       "to_do",
		},
		{
			UserId:      "0002",
			Description: "description 2",
			State:       "to_do",
		},
	}

	rowTasks := sqlmock.NewRows([]string{"id", "user_id", "description", "state"}).
		AddRow("", listTasks[0].UserId, listTasks[0].Description, listTasks[0].State).
		AddRow("", listTasks[1].UserId, listTasks[1].Description, listTasks[1].State)

	tests := map[string]struct {
		dbError      error
		dbRowTask    *sqlmock.Rows
		expectedResp []common.Task
		expectedErr  error
	}{
		"success": {
			dbError:      nil,
			dbRowTask:    rowTasks,
			expectedResp: listTasks,
			expectedErr:  nil,
		},
		"fail": {
			dbError:      errGetTask,
			dbRowTask:    nil,
			expectedResp: nil,
			expectedErr:  errGetTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			mockGet := mock.ExpectQuery("SELECT id, user_id, description, state FROM public.task")
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowTask)
			} else {
				mockGet.WillReturnError(errGetTask)
			}

			task, err := db.ListTasks()
			assert.Equal(t, task, test.expectedResp)
			assert.Equal(t, err, test.expectedErr)
		})

	}
}

func TestUserTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := Config{
		Logger: logger,
	}

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errGetTask := errors.New("any error")
	listTasks := []common.Task{
		{
			UserId:      "0001",
			Description: "description 1",
			State:       "to_do",
		},
		{
			UserId:      "0001",
			Description: "description 2",
			State:       "to_do",
		},
	}

	rowTasks := sqlmock.NewRows([]string{"id", "user_id", "description", "state"}).
		AddRow("", listTasks[0].UserId, listTasks[0].Description, listTasks[0].State).
		AddRow("", listTasks[1].UserId, listTasks[1].Description, listTasks[1].State)

	tests := map[string]struct {
		id           string
		dbError      error
		dbRowTask    *sqlmock.Rows
		expectedResp []common.Task
		expectedErr  error
	}{
		"success": {
			dbError:      nil,
			dbRowTask:    rowTasks,
			expectedResp: listTasks,
			expectedErr:  nil,
		},
		"fail": {
			dbError:      errGetTask,
			dbRowTask:    nil,
			expectedResp: nil,
			expectedErr:  errGetTask,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			mockGet := mock.ExpectQuery("SELECT id, user_id, description, state FROM public.task")
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowTask)
			} else {
				mockGet.WillReturnError(errGetTask)
			}

			task, err := db.ListUserTasks(test.id)
			assert.Equal(t, task, test.expectedResp)
			assert.Equal(t, err, test.expectedErr)
		})

	}
}
