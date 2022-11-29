package db

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
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
			mockInsert := mock.ExpectExec("INSERT INTO public.user").WithArgs(test.user.Id, test.user.Username, test.user.Name)
			if test.dbError == nil {
				mockInsert.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockInsert.WillReturnError(test.dbError)
			}

			err := db.AddUser(test.user)
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

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errAddUser := errors.New("error updating user")
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
			mockUpdate := mock.ExpectExec("UPDATE public.user").WithArgs(test.user.Username, test.user.Name, test.user.Id)
			if test.dbError == nil {
				mockUpdate.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockUpdate.WillReturnError(test.dbError)
			}

			err := db.UpdateUser(test.user)
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

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errGetUser := errors.New("any error")
	user := &common.User{
		Username: "username1",
		Name:     "User Name 1",
	}
	rowUser := sqlmock.NewRows([]string{"id", "username", "name"}).AddRow("", user.Username, user.Name)

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
		t.Run(index, func(t *testing.T) {
			mockGet := mock.ExpectQuery("SELECT id, username, name FROM public.user").WithArgs(test.id)
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowUser)
			} else {
				mockGet.WillReturnError(errGetUser)
			}

			user, err := db.GetUser(test.id)
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

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errAddUser := errors.New("error deleting user")

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
			dbError:      errAddUser,
			expectedResp: errAddUser,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			mockUpdate := mock.ExpectExec("DELETE FROM public.user").WithArgs(test.id)
			if test.dbError == nil {
				mockUpdate.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockUpdate.WillReturnError(test.dbError)
			}

			err := db.DeleteUser(test.id)
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

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	db, err := New(cfg)
	db.db = dbMock

	assert.NoError(t, err)

	errGetUser := errors.New("any error")
	listUsers := []common.User{
		{
			Username: "username1",
			Name:     "User Name 1",
		},
		{
			Username: "username2",
			Name:     "User Name 2",
		},
	}

	rowUsers := sqlmock.NewRows([]string{"id", "username", "name"}).
		AddRow("", listUsers[0].Username, listUsers[0].Name).
		AddRow("", listUsers[1].Username, listUsers[1].Name)

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
			dbError:      errGetUser,
			dbRowUser:    nil,
			expectedResp: nil,
			expectedErr:  errGetUser,
		},
	}

	for index, test := range tests {
		t.Run(index, func(t *testing.T) {
			mockGet := mock.ExpectQuery("SELECT id, username, name FROM public.user")
			if test.dbError == nil {
				mockGet.WillReturnRows(test.dbRowUser)
			} else {
				mockGet.WillReturnError(errGetUser)
			}

			user, err := db.ListUsers()
			assert.Equal(t, user, test.expectedResp)
			assert.Equal(t, err, test.expectedErr)
		})

	}
}
