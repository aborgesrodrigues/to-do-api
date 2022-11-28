package service

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (svc *Service) AddTask(w http.ResponseWriter, r *http.Request) {
	request := &common.Task{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		svc.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// add uuid
	request.Id = uuid.New().String()

	if err := svc.DB.AddTask(request); err != nil {
		svc.Logger.Error("Unable add Task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, map[string]string{
		"message": "Task Added",
	})
}

func (svc *Service) UpdateTask(w http.ResponseWriter, r *http.Request) {
	request := &common.Task{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		svc.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	request.Id = r.Context().Value(taskIdCtx).(string)

	if err := svc.DB.UpdateTask(request); err != nil {
		svc.Logger.Error("Unable add Task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, map[string]string{
		"message": "Task Updated",
	})
}

func (svc *Service) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(taskIdCtx).(string)

	task, err := svc.DB.GetTask(id)
	if err != nil {
		svc.Logger.Error("Unable to retrieve task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// get user
	user, err := svc.DB.GetUser(task.UserId)
	if err != nil {
		svc.Logger.Error("Unable to retrieve task user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	task.User = user

	writeResponse(w, http.StatusOK, task)
}

func (svc *Service) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(taskIdCtx).(string)

	err := svc.DB.DeleteTask(id)
	if err != nil {
		svc.Logger.Error("Unable to delete tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, nil)
}

func (svc *Service) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := svc.DB.ListTasks()
	if err != nil {
		svc.Logger.Error("Unable to retrieve tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, tasks)
}
