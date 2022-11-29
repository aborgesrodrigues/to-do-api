package service

import (
	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (svc *Service) AddUser(user *common.User) (*common.User, error) {
	user.Id = uuid.New().String()

	if err := svc.db.AddUser(user); err != nil {
		svc.logger.Error("Unable add user.", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (svc *Service) UpdateUser(user *common.User) (*common.User, error) {
	if err := svc.db.UpdateUser(user); err != nil {
		svc.logger.Error("Unable add user.", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (svc *Service) GetUser(id string) (*common.User, error) {
	user, err := svc.db.GetUser(id)
	if err != nil {
		svc.logger.Error("Unable to retrieve users.", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (svc *Service) DeleteUser(id string) error {
	// delete user tasks
	if err := svc.db.DeleteUserTasks(id); err != nil {
		svc.logger.Error("Unable to delete user tasks.", zap.Error(err))
		return err
	}

	// delete user
	if err := svc.db.DeleteUser(id); err != nil {
		svc.logger.Error("Unable to delete users.", zap.Error(err))
		return err
	}

	return nil
}

func (svc *Service) ListUsers() ([]common.User, error) {
	users, err := svc.db.ListUsers()
	if err != nil {
		svc.logger.Error("Unable to retrieve users.", zap.Error(err))
		return nil, err
	}

	return users, nil
}

func (svc *Service) ListUserTasks(id string) ([]common.Task, error) {
	users, err := svc.db.ListUserTasks(id)
	if err != nil {
		svc.logger.Error("Unable to retrieve user tasks.", zap.Error(err))
		return nil, err
	}

	return users, nil
}
