package db

import (
	"errors"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/aborgesrodrigues/to-do-api/internal/common"
)

func (d *dbTestSuite) TestAddTask() {
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
		d.Run(index, func() {
			mockInsert := d.mock.ExpectExec("INSERT INTO public.task").WithArgs(test.task.Id, test.task.UserId, test.task.Description, test.task.State)
			if test.dbError == nil {
				mockInsert.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockInsert.WillReturnError(test.dbError)
			}

			err := d.db.AddTask(test.task)
			d.Assert().Equal(err, test.expectedResp)
		})

	}
}

func (d *dbTestSuite) TestUpdateTask() {
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
		d.Run(index, func() {
			mockUpdate := d.mock.ExpectExec("UPDATE public.task").WithArgs(test.task.UserId, test.task.Description, test.task.State, test.task.Id)
			if test.dbError == nil {
				mockUpdate.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockUpdate.WillReturnError(test.dbError)
			}

			err := d.db.UpdateTask(test.task)
			d.Assert().Equal(err, test.expectedResp)
		})

	}
}

func (d *dbTestSuite) TestGetTask() {
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
		d.Run(index, func() {
			mockGet := d.mock.ExpectQuery("SELECT id, user_id, description, state FROM public.task").WithArgs(test.id)
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowTask)
			} else {
				mockGet.WillReturnError(errGetTask)
			}

			task, err := d.db.GetTask(test.id)
			d.Assert().Equal(task, test.expectedResp)
			d.Assert().Equal(err, test.expectedErr)
		})

	}
}

func (d *dbTestSuite) TestDeleteTask() {
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
		d.Run(index, func() {
			mockDelete := d.mock.ExpectExec("DELETE FROM public.task").WithArgs(test.id)
			if test.dbError == nil {
				mockDelete.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockDelete.WillReturnError(test.dbError)
			}

			err := d.db.DeleteTask(test.id)
			d.Assert().Equal(err, test.expectedResp)
		})

	}
}

func (d *dbTestSuite) TestDeleteUserTasks() {
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
		d.Run(index, func() {
			mockDelete := d.mock.ExpectExec("DELETE FROM public.task").WithArgs(test.id)
			if test.dbError == nil {
				mockDelete.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockDelete.WillReturnError(test.dbError)
			}

			err := d.db.DeleteUserTasks(test.id)
			d.Assert().Equal(err, test.expectedResp)
		})

	}
}

func (d *dbTestSuite) TestListTasks() {
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
		d.Run(index, func() {
			mockGet := d.mock.ExpectQuery("SELECT id, user_id, description, state FROM public.task")
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowTask)
			} else {
				mockGet.WillReturnError(errGetTask)
			}

			task, err := d.db.ListTasks()
			d.Assert().Equal(task, test.expectedResp)
			d.Assert().Equal(err, test.expectedErr)
		})

	}
}

func (d *dbTestSuite) TestUserTasks() {
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
		d.Run(index, func() {
			mockGet := d.mock.ExpectQuery("SELECT id, user_id, description, state FROM public.task")
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowTask)
			} else {
				mockGet.WillReturnError(errGetTask)
			}

			task, err := d.db.ListUserTasks(test.id)
			d.Assert().Equal(task, test.expectedResp)
			d.Assert().Equal(err, test.expectedErr)
		})

	}
}
