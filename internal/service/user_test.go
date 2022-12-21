package service

import (
	"errors"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
)

func (s *svcTestSuite) TestAddUser() {
	errAddUser := errors.New("error inserting user")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}

	tests := map[string]struct {
		user        *common.User
		dbError     error
		expectedErr error
	}{
		"success": {
			user:        user,
			dbError:     nil,
			expectedErr: nil,
		},
		"fail": {
			user:        user,
			dbError:     errAddUser,
			expectedErr: errAddUser,
		},
	}

	for index, test := range tests {
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				AddUser(test.user).
				Return(test.dbError)

			user, err := s.svc.AddUser(test.user)
			s.Assert().Equal(err, test.expectedErr)

			if test.dbError == nil {
				s.Assert().NotEmpty(user.Id)
			}
		})

	}
}

func (s *svcTestSuite) TestUpdateUser() {
	errAddUser := errors.New("error inserting user")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}

	tests := map[string]struct {
		user        *common.User
		dbError     error
		expectedErr error
	}{
		"success": {
			user:        user,
			dbError:     nil,
			expectedErr: nil,
		},
		"fail": {
			user:        user,
			dbError:     errAddUser,
			expectedErr: errAddUser,
		},
	}

	for index, test := range tests {
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				UpdateUser(test.user).
				Return(test.dbError)

			_, err := s.svc.UpdateUser(test.user)
			s.Assert().Equal(err, test.expectedErr)
		})

	}
}

func (s *svcTestSuite) TestGetUser() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				GetUser(test.id).
				Return(test.dbUser, test.dbError)

			user, err := s.svc.GetUser(test.id)
			s.Assert().Equal(user, test.expectedResp)
			s.Assert().Equal(err, test.expectedErr)
		})

	}
}

func (s *svcTestSuite) TestDeleteUser() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				DeleteUserTasks(test.id).
				Return(test.dbError1)

			if test.dbError1 == nil {
				s.getDB().
					DeleteUser(test.id).
					Return(test.dbError2)
			}

			err := s.svc.DeleteUser(test.id)
			s.Assert().Equal(err, test.expectedResp)
		})
	}
}

func (s *svcTestSuite) TestListUsers() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				ListUsers().
				Return(test.dbUsers, test.dbError)

			user, err := s.svc.ListUsers()
			s.Assert().Equal(user, test.expectedResp)
			s.Assert().Equal(err, test.expectedErr)
		})

	}
}

func (s *svcTestSuite) TestListUserTasks() {
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
		s.Run(index, func() {
			// set up dao mock
			s.getDB().
				ListUserTasks(test.id).
				Return(test.dbTasks, test.dbError)

			users, err := s.svc.ListUserTasks(test.id)
			s.Assert().Equal(users, test.expectedResp)
			s.Assert().Equal(err, test.expectedErr)
		})

	}
}
