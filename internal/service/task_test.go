package service

import (
	"errors"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
)

func (s *svcTestSuite) TestAddTask() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				AddTask(test.task).
				Return(test.dbError)

			task, err := s.svc.AddTask(test.task)
			s.Assert().Equal(err, test.expectedResp)

			if test.dbError == nil {
				s.Assert().NotEmpty(task.Id)
			}
		})

	}
}

func (s *svcTestSuite) TestUpdateTask() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				UpdateTask(test.task).
				Return(test.dbError)

			_, err := s.svc.UpdateTask(test.task)
			s.Assert().Equal(err, test.expectedResp)
		})

	}
}

func (s *svcTestSuite) TestGetTask() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				GetTask(test.id).
				Return(test.dbTask, test.dbError1)

			if test.dbError1 == nil {
				s.getDB().
					GetUser("00001").
					Return(&common.User{}, test.dbError2)
			}

			task, err := s.svc.GetTask(test.id)
			s.Assert().Equal(task, test.expectedResp)
			s.Assert().Equal(err, test.expectedErr)
		})

	}
}

func (s *svcTestSuite) TestDeleteTask() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				DeleteTask(test.id).
				Return(test.dbError)

			err := s.svc.DeleteTask(test.id)
			s.Assert().Equal(err, test.expectedResp)
		})
	}
}

func (s *svcTestSuite) TestListTasks() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				ListTasks().
				Return(test.dbTasks, test.dbError)

			tasks, err := s.svc.ListTasks()
			s.Assert().Equal(tasks, test.expectedResp)
			s.Assert().Equal(err, test.expectedErr)
		})

	}
}
