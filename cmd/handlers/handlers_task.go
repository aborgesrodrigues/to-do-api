package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
)

// UpdateTask godoc
// @Summary      Update a task
// @Description  Updates the task description and state. Only the task owner can update it.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Task ID"
// @Param        task  body      common.Task  true  "Task data"
// @Success      200   {object}  common.Task
// @Failure      403   {string}  string  "Task does not belong to the authenticated user"
// @Failure      500   {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /tasks/{id} [put]
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

// GetTask godoc
// @Summary      Get a task by ID
// @Description  Only the task owner can retrieve it.
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  common.Task
// @Failure      403  {string}  string  "Task does not belong to the authenticated user"
// @Failure      500  {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /tasks/{id} [get]
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

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Only the task owner can delete it.
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  common.MessageResponse
// @Failure      403  {string}  string  "Task does not belong to the authenticated user"
// @Failure      500  {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /tasks/{id} [delete]
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

// ListTasks godoc
// @Summary      List all tasks
// @Tags         tasks
// @Produce      json
// @Success      200  {array}   common.Task
// @Failure      500  {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /tasks [get]
func (handler *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := handler.svc.ListTasks()
	if err != nil {
		handler.Logger.Error("Unable to retrieve tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, tasks)
}
