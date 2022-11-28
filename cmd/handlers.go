package main

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
)

func (handler *handler) addUser(w http.ResponseWriter, r *http.Request) {
	request := &common.User{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := handler.svc.AddUser(request); err != nil {
		handler.logger.Error("Unable add user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, map[string]string{
		"message": "User Added",
	})
}

func (handler *handler) updateUser(w http.ResponseWriter, r *http.Request) {
	request := &common.User{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	request.Id = r.Context().Value(userIdCtx).(string)

	if err := handler.svc.UpdateUser(request); err != nil {
		handler.logger.Error("Unable add user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, map[string]string{
		"message": "User Updated",
	})
}

func (handler *handler) getUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIdCtx).(string)

	user, err := handler.svc.GetUser(id)
	if err != nil {
		handler.logger.Error("Unable to retrieve users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, user)
}

func (handler *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIdCtx).(string)

	// delete user
	if err := handler.svc.DeleteUser(id); err != nil {
		handler.logger.Error("Unable to delete users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, map[string]string{
		"message": "User Deleted",
	})
}

func (handler *handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.svc.ListUsers()
	if err != nil {
		handler.logger.Error("Unable to retrieve users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, users)
}

func (handler *handler) getUserTasks(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIdCtx).(string)

	users, err := handler.svc.ListUserTasks(id)
	if err != nil {
		handler.logger.Error("Unable to retrieve user tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, users)
}

func (handler *handler) addTask(w http.ResponseWriter, r *http.Request) {
	request := &common.Task{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := handler.svc.AddTask(request); err != nil {
		handler.logger.Error("Unable add Task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, map[string]string{
		"message": "Task Added",
	})
}

func (handler *handler) updateTask(w http.ResponseWriter, r *http.Request) {
	request := &common.Task{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.logger.Error("Unable to decode request body.", zap.Error(err))
		return
	}

	request.Id = r.Context().Value(taskIdCtx).(string)

	if err := handler.svc.UpdateTask(request); err != nil {
		handler.logger.Error("Unable add Task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, map[string]string{
		"message": "Task Updated",
	})
}

func (handler *handler) getTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(taskIdCtx).(string)

	task, err := handler.svc.GetTask(id)
	if err != nil {
		handler.logger.Error("Unable to retrieve task.", zap.Error(err))
		return
	}

	writeResponse(w, http.StatusOK, task)
}

func (handler *handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(taskIdCtx).(string)
	err := handler.svc.DeleteTask(id)
	if err != nil {
		handler.logger.Error("Unable to delete tasks.", zap.Error(err))
		return
	}

	writeResponse(w, http.StatusOK, map[string]string{
		"message": "Task deleted",
	})
}

func (handler *handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := handler.svc.ListTasks()
	if err != nil {
		handler.logger.Error("Unable to retrieve tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, tasks)
}
