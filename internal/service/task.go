package service

import (
	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (svc *Service) AddTask(task *common.Task) error {
	// add uuid
	task.Id = uuid.New().String()

	if err := svc.db.AddTask(task); err != nil {
		svc.logger.Error("Unable add Task.", zap.Error(err))
		return err
	}

	return nil
}

func (svc *Service) UpdateTask(task *common.Task) error {
	if err := svc.db.UpdateTask(task); err != nil {
		svc.logger.Error("Unable add Task.", zap.Error(err))
		return err
	}

	return nil
}

func (svc *Service) GetTask(id string) (*common.Task, error) {
	task, err := svc.db.GetTask(id)
	if err != nil {
		svc.logger.Error("Unable to retrieve task.", zap.Error(err))
		return nil, err
	}

	// get user
	user, err := svc.db.GetUser(task.UserId)
	if err != nil {
		svc.logger.Error("Unable to retrieve task user.", zap.Error(err))
		return nil, err
	}
	task.User = user

	return task, nil
}

func (svc *Service) DeleteTask(id string) error {
	err := svc.db.DeleteTask(id)
	if err != nil {
		svc.logger.Error("Unable to delete tasks.", zap.Error(err))
		return err
	}

	return nil
}

func (svc *Service) ListTasks() ([]common.Task, error) {
	tasks, err := svc.db.ListTasks()
	if err != nil {
		svc.logger.Error("Unable to retrieve tasks.", zap.Error(err))
		return nil, err
	}

	return tasks, nil
}
