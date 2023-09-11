package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
)

func (handler *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	request := &common.User{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := handler.svc.AddUser(request)
	if err != nil {
		handler.Logger.Error("Unable add user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := generateJWT(user)
	if err != nil {
		handler.Logger.Error("Error generating JWT.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated,
		map[string]any{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	request := &common.User{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := handler.svc.Login(request.Username, request.Password)
	if err != nil {
		handler.Logger.Error("Unable to login.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := generateJWT(user)
	if err != nil {
		handler.Logger.Error("Error generating JWT.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK,
		map[string]any{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
}

func (handler *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	request := &common.User{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	request.Id = r.Context().Value(idCtx).(string)

	user, err := handler.svc.UpdateUser(request)
	if err != nil {
		handler.Logger.Error("Unable update user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, user)
}

func (handler *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idCtx).(string)

	user, err := handler.svc.GetUser(id)
	if err != nil {
		handler.Logger.Error("Unable to retrieve user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, user)
}

func (handler *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idCtx).(string)

	// delete user
	if err := handler.svc.DeleteUser(id); err != nil {
		handler.Logger.Error("Unable to delete users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, map[string]string{
		"message": "User Deleted",
	})
}

func (handler *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.svc.ListUsers()
	if err != nil {
		handler.Logger.Error("Unable to retrieve users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, users)
}

func (handler *Handler) ListUserTasks(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idCtx).(string)

	users, err := handler.svc.ListUserTasks(id)
	if err != nil {
		handler.Logger.Error("Unable to retrieve user tasks.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, users)
}

func (handler *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idCtx).(string)

	user, err := handler.svc.GetUser(id)
	if err != nil {
		handler.Logger.Error("Unable to retrieve user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := generateJWT(user)
	if err != nil {
		handler.Logger.Error("Error generating JWT.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated,
		map[string]any{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
}
