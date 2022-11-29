package service

import (
	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/aborgesrodrigues/to-do-api/internal/db"
	"go.uber.org/zap"
)

type SVCInterface interface {
	AddTask(task *common.Task) (*common.Task, error)
	UpdateTask(task *common.Task) (*common.Task, error)
	GetTask(id string) (*common.Task, error)
	DeleteTask(id string) error
	ListTasks() ([]common.Task, error)
	ListUserTasks(id string) ([]common.Task, error)

	AddUser(user *common.User) (*common.User, error)
	UpdateUser(user *common.User) (*common.User, error)
	GetUser(id string) (*common.User, error)
	DeleteUser(id string) error
	ListUsers() ([]common.User, error)
}

type Config struct {
	Logger *zap.Logger
}

type Service struct {
	logger *zap.Logger
	db     db.DBInterface
}
