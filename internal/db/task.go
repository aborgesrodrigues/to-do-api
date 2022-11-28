package db

import (
	"github.com/aborgesrodrigues/to-do-api/internal/common"
)

func (db *DB) AddTask(task *common.Task) error {
	_, err := db.db.Exec(`
		INSERT INTO public.task(id, user_id, description, state)
		VALUES($1, $2, $3, $4)
	`, task.Id, task.UserId, task.Description, task.State)

	if err != nil {
		db.logger.Error("Error inserting task.")
		return err
	}

	return nil
}

func (db *DB) UpdateTask(task *common.Task) error {
	_, err := db.db.Exec(`
		UPDATE public.task
		SET user_id = $1, description = $2, state = $3
		WHERE id = $4
	`, task.UserId, task.Description, task.State, task.Id)

	if err != nil {
		db.logger.Error("Error updating task.")
		return err
	}

	return nil
}

func (db *DB) GetTask(id string) (*common.Task, error) {
	results, err := db.db.Query(`
		SELECT id, user_id, description, state 
		FROM public.task
		WHERE id= $1`, id)

	if err != nil {
		db.logger.Error("Error retrieving task.")
		return nil, err
	}

	task := common.Task{}
	for results.Next() {
		err = results.Scan(
			&task.Id,
			&task.UserId,
			&task.Description,
			&task.State)
		if err != nil {
			db.logger.Error("Error mapping database data to struct.")
			return nil, err
		}
	}
	return &task, nil
}

func (db *DB) DeleteTask(id string) error {
	_, err := db.db.Exec(`
		DELETE FROM public.task WHERE id = $1
	`, id)

	if err != nil {
		db.logger.Error("Error deleting task.")
		return err
	}

	return nil
}

func (db *DB) DeleteUserTasks(userId string) error {
	_, err := db.db.Exec(`
		DELETE FROM public.task WHERE user_id = $1
	`, userId)

	if err != nil {
		db.logger.Error("Error deleting user tasks.")
		return err
	}

	return nil
}

func (db *DB) ListTasks() ([]common.Task, error) {
	results, err := db.db.Query(`
		SELECT id, user_id, description, state 
		FROM public.task`)

	if err != nil {
		return nil, err
	}

	tasks := make([]common.Task, 0)
	for results.Next() {
		task := common.Task{}
		err = results.Scan(
			&task.Id,
			&task.UserId,
			&task.Description,
			&task.State)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *DB) ListUserTasks(id string) ([]common.Task, error) {
	results, err := db.db.Query(`
		SELECT id, user_id, description, state 
		FROM public.task
		WHERE user_id = $1`, id)

	if err != nil {
		return nil, err
	}

	tasks := make([]common.Task, 0)
	for results.Next() {
		task := common.Task{}
		err = results.Scan(
			&task.Id,
			&task.UserId,
			&task.Description,
			&task.State)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
