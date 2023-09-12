package db

import (
	"errors"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
)

var errLogin = errors.New("invalid username/password")

func (db *DB) AddUser(user *common.User) error {
	_, err := db.db.Exec(`
		INSERT INTO public.user(id, username, name, password)
		VALUES($1, $2, $3, $4)
	`, user.Id, user.Username, user.Name, user.Password)

	if err != nil {
		db.logger.Error("Error inserting user.")
		return err
	}

	return nil
}

func (db *DB) UpdateUser(user *common.User) error {
	_, err := db.db.Exec(`
		UPDATE public.user
		SET username = $1, name = $2
		WHERE id = $3
	`, user.Username, user.Name, user.Id)

	if err != nil {
		db.logger.Error("Error updating user.")
		return err
	}

	return nil
}

func (db *DB) GetUser(id string) (*common.User, error) {
	results, err := db.db.Query(`
		SELECT id, username, name, password
		FROM public.user
		WHERE id= $1`, id)

	if err != nil {
		db.logger.Error("Error retrieving user.")
		return nil, err
	}

	user := common.User{}
	for results.Next() {
		err = results.Scan(
			&user.Id,
			&user.Username,
			&user.Name,
			&user.Password)
		if err != nil {
			db.logger.Error("Error mapping database data to struct.")
			return nil, err
		}
	}
	return &user, nil
}

func (db *DB) GetUserByUsernamePassword(username, password string) (*common.User, error) {
	results, err := db.db.Query(`
		SELECT id, username, name, password
		FROM public.user
		WHERE username= $1`, username)

	if err != nil {
		db.logger.Error("Error retrieving user.")
		return nil, err
	}

	user := common.User{}
	for results.Next() {
		err = results.Scan(
			&user.Id,
			&user.Username,
			&user.Name,
			&user.Password)
		if err != nil {
			db.logger.Error("Error mapping database data to struct.")
			return nil, err
		}

		if !common.CheckPasswordHash(password, user.Password) {
			db.logger.Error("Invalid password.")
			return nil, errLogin
		}

	}

	if user.Id == "" {
		db.logger.Error("User not found.")
		return nil, errLogin
	}

	return &user, nil
}

func (db *DB) DeleteUser(id string) error {
	_, err := db.db.Exec(`
		DELETE FROM public.user WHERE id = $1
	`, id)

	if err != nil {
		db.logger.Error("Error deleting user.")
		return err
	}

	return nil
}

func (db *DB) ListUsers() ([]common.User, error) {
	results, err := db.db.Query(`
		SELECT id, username, name, password
		FROM public.user`)

	if err != nil {
		db.logger.Error("Error retrieving users.")
		return nil, err
	}

	users := make([]common.User, 0)
	for results.Next() {
		user := common.User{}
		err = results.Scan(
			&user.Id,
			&user.Username,
			&user.Name,
			&user.Password)
		if err != nil {
			db.logger.Error("Error mapping database data to struct.")
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
