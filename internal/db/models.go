package db

import (
	"database/sql"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
)

type DBInterface interface {
	AddTask(task *common.Task) error
	UpdateTask(task *common.Task) error
	GetTask(id string) (*common.Task, error)
	DeleteTask(id string) error
	DeleteUserTasks(userId string) error
	ListTasks() ([]common.Task, error)
	ListUserTasks(id string) ([]common.Task, error)

	AddUser(user *common.User) error
	UpdateUser(user *common.User) error
	GetUser(id string) (*common.User, error)
	DeleteUser(id string) error
	ListUsers() ([]common.User, error)
}

type Config struct {
	Logger *zap.Logger
}

type DB struct {
	db     *sql.DB
	logger *zap.Logger
}
