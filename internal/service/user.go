package service

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (svc *Service) AddUser(w http.ResponseWriter, r *http.Request) {
	request := &common.User{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		svc.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// add uuid
	request.Id = uuid.New().String()

	if err := svc.DB.AddUser(request); err != nil {
		svc.Logger.Error("Unable add user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, map[string]string{
		"message": "User Added",
	})
}

func (svc *Service) UpdateUser(w http.ResponseWriter, r *http.Request) {
	request := &common.User{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		svc.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	request.Id = r.Context().Value(userIdCtx).(string)

	if err := svc.DB.UpdateUser(request); err != nil {
		svc.Logger.Error("Unable add user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, map[string]string{
		"message": "User Updated",
	})
}

func (svc *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIdCtx).(string)

	users, err := svc.DB.GetUser(id)
	if err != nil {
		svc.Logger.Error("Unable to retrieve users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, users)
}

func (svc *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIdCtx).(string)

	// delete user tasks
	if err := svc.DB.DeleteTasksUser(id); err != nil {
		svc.Logger.Error("Unable to delete user tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// delete user
	if err := svc.DB.DeleteUser(id); err != nil {
		svc.Logger.Error("Unable to delete users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, map[string]string{
		"message": "User Deleted",
	})
}

func (svc *Service) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := svc.DB.ListUsers()
	if err != nil {
		svc.Logger.Error("Unable to retrieve users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, users)
}

func (svc *Service) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIdCtx).(string)

	users, err := svc.DB.ListUserTasks(id)
	if err != nil {
		svc.Logger.Error("Unable to retrieve user tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, users)
}
