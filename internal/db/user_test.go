package db

import (
	"errors"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/aborgesrodrigues/to-do-api/internal/common"
)

func (d *dbTestSuite) TestAddUser() {
	errAddUser := errors.New("error inserting user")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
		Password: "password",
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
		d.Run(index, func() {
			mockInsert := d.mock.ExpectExec("INSERT INTO public.user").WithArgs(test.user.Id, test.user.Username, test.user.Name, test.user.Password)
			if test.dbError == nil {
				mockInsert.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockInsert.WillReturnError(test.dbError)
			}

			err := d.db.AddUser(test.user)
			d.Assert().Equal(err, test.expectedResp)
		})

	}
}

func (d *dbTestSuite) TestUpdateUser() {
	errUpdateUser := errors.New("error updating user")
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
			dbError:      errUpdateUser,
			expectedResp: errUpdateUser,
		},
	}

	for index, test := range tests {
		d.Run(index, func() {
			mockUpdate := d.mock.ExpectExec("UPDATE public.user").WithArgs(test.user.Username, test.user.Name, test.user.Id)
			if test.dbError == nil {
				mockUpdate.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockUpdate.WillReturnError(test.dbError)
			}

			err := d.db.UpdateUser(test.user)
			d.Assert().Equal(err, test.expectedResp)
		})

	}
}

func (d *dbTestSuite) TestGetUser() {
	errGetUser := errors.New("any error")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
		Password: "password",
	}
	rowUser := sqlmock.NewRows([]string{"id", "username", "name", "password"}).AddRow("", user.Username, user.Name, user.Password)

	tests := map[string]struct {
		id           string
		dbError      error
		dbRowUser    *sqlmock.Rows
		expectedResp *common.User
		expectedErr  error
	}{
		"success": {
			dbError:      nil,
			dbRowUser:    rowUser,
			expectedResp: user,
			expectedErr:  nil,
		},
		"fail": {
			dbError:      errGetUser,
			dbRowUser:    nil,
			expectedResp: nil,
			expectedErr:  errGetUser,
		},
	}

	for index, test := range tests {
		d.Run(index, func() {
			mockGet := d.mock.ExpectQuery("SELECT id, username, name, password FROM public.user").WithArgs(test.id)
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowUser)
			} else {
				mockGet.WillReturnError(errGetUser)
			}

			user, err := d.db.GetUser(test.id)
			d.Assert().Equal(user, test.expectedResp)
			d.Assert().Equal(err, test.expectedErr)
		})

	}
}

func (d *dbTestSuite) TestDeleteUser() {
	errDeleteUser := errors.New("error deleting user")

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
			dbError:      errDeleteUser,
			expectedResp: errDeleteUser,
		},
	}

	for index, test := range tests {
		d.Run(index, func() {
			mockUpdate := d.mock.ExpectExec("DELETE FROM public.user").WithArgs(test.id)
			if test.dbError == nil {
				mockUpdate.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockUpdate.WillReturnError(test.dbError)
			}

			err := d.db.DeleteUser(test.id)
			d.Assert().Equal(err, test.expectedResp)
		})

	}
}

func (d *dbTestSuite) TestListUsers() {
	errListUsers := errors.New("any error")
	listUsers := []common.User{
		{
			Username: "username1",
			Name:     "User Name 1",
			Password: "password1",
		},
		{
			Username: "username2",
			Name:     "User Name 2",
			Password: "password2",
		},
	}

	rowUsers := sqlmock.NewRows([]string{"id", "username", "name", "password"}).
		AddRow("", listUsers[0].Username, listUsers[0].Name, listUsers[0].Password).
		AddRow("", listUsers[1].Username, listUsers[1].Name, listUsers[1].Password)

	tests := map[string]struct {
		dbError      error
		dbRowUser    *sqlmock.Rows
		expectedResp []common.User
		expectedErr  error
	}{
		"success": {
			dbError:      nil,
			dbRowUser:    rowUsers,
			expectedResp: listUsers,
			expectedErr:  nil,
		},
		"fail": {
			dbError:      errListUsers,
			dbRowUser:    nil,
			expectedResp: nil,
			expectedErr:  errListUsers,
		},
	}

	for index, test := range tests {
		d.Run(index, func() {
			mockGet := d.mock.ExpectQuery("SELECT id, username, name, password FROM public.user")
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowUser)
			} else {
				mockGet.WillReturnError(errListUsers)
			}

			user, err := d.db.ListUsers()
			d.Assert().Equal(user, test.expectedResp)
			d.Assert().Equal(err, test.expectedErr)
		})

	}
}
