package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
)

func (handler *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(idCtx).(string)

	request := &common.Task{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	request.User = &common.User{
		Id: userID,
	}
	task, err := handler.svc.AddTask(request)
	if err != nil {
		handler.Logger.Error("Unable add Task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, task)
}

func (handler *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	request := &common.Task{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.Logger.Error("Unable to decode request body.", zap.Error(err))
		return
	}

	savedTask, err := handler.svc.GetTask(r.Context().Value(idCtx).(string))
	if err != nil {
		handler.Logger.Error("Unable retrieve task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// check if task belong to user
	claims := r.Context().Value(claimsKey).(*common.Claims)
	if savedTask.User.Id != claims.UserID {
		handler.Logger.Error("Id doesn't match with token.")
		writeResponse(w, http.StatusForbidden, "Id doesn't match with token.")
		return
	}

	savedTask.Description = request.Description
	savedTask.State = request.State

	task, err := handler.svc.UpdateTask(savedTask)
	if err != nil {
		handler.Logger.Error("Unable add Task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, task)
}

func (handler *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idCtx).(string)

	task, err := handler.svc.GetTask(id)
	if err != nil {
		handler.Logger.Error("Unable to retrieve task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// check if task belong to user
	claims := r.Context().Value(claimsKey).(*common.Claims)
	if task.User.Id != claims.UserID {
		handler.Logger.Error("Id doesn't match with token.")
		writeResponse(w, http.StatusForbidden, "Id doesn't match with token.")
		return
	}

	writeResponse(w, http.StatusOK, task)
}

func (handler *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idCtx).(string)

	task, err := handler.svc.GetTask(id)
	if err != nil {
		handler.Logger.Error("Unable to retrieve task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// check if task belong to user
	claims := r.Context().Value(claimsKey).(*common.Claims)
	if task.User.Id != claims.UserID {
		handler.Logger.Error("Id doesn't match with token.")
		writeResponse(w, http.StatusForbidden, "Id doesn't match with token.")
		return
	}

	err = handler.svc.DeleteTask(id)
	if err != nil {
		handler.Logger.Error("Unable to delete tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, map[string]string{
		"message": "Task Deleted",
	})
}

func (handler *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := handler.svc.ListTasks()
	if err != nil {
		handler.Logger.Error("Unable to retrieve tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, tasks)
}
