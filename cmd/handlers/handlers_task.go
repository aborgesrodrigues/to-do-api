package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
)

func (handler *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	request := &common.Task{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	task, err := handler.svc.AddTask(request)
	if err != nil {
		handler.logger.Error("Unable add Task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, task)
}

func (handler *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	request := &common.Task{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.logger.Error("Unable to decode request body.", zap.Error(err))
		return
	}

	request.Id = r.Context().Value(taskIdCtx).(string)

	task, err := handler.svc.UpdateTask(request)
	if err != nil {
		handler.logger.Error("Unable add Task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, task)
}

func (handler *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(taskIdCtx).(string)

	task, err := handler.svc.GetTask(id)
	if err != nil {
		handler.logger.Error("Unable to retrieve task.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, task)
}

func (handler *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(taskIdCtx).(string)
	err := handler.svc.DeleteTask(id)
	if err != nil {
		handler.logger.Error("Unable to delete tasks.", zap.Error(err))
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
		handler.logger.Error("Unable to retrieve tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, tasks)
}
