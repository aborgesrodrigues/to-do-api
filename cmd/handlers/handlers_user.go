package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
)

// AddUser godoc
// @Summary      Create a user
// @Description  Creates a user and returns the user data along with an access token and a refresh token.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      common.User  true  "User data"
// @Success      201   {object}  common.AuthResponse
// @Failure      500   {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /users [post]
// @Router       /token [post]
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

// UpdateUser godoc
// @Summary      Update a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "User ID"
// @Param        user  body      common.User  true  "User data"
// @Success      200   {object}  common.User
// @Failure      500   {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{id} [put]
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

// GetUser godoc
// @Summary      Get a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  common.User
// @Failure      500  {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{id} [get]
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

// DeleteUser godoc
// @Summary      Delete a user
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  common.MessageResponse
// @Failure      500  {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{id} [delete]
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

// ListUsers godoc
// @Summary      List users
// @Tags         users
// @Produce      json
// @Success      200  {array}   common.User
// @Failure      500  {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /users [get]
func (handler *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.svc.ListUsers()
	if err != nil {
		handler.Logger.Error("Unable to retrieve users.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Logger.Info("Users listed.",
		zap.Bool("success", true))

	writeResponse(w, http.StatusOK, users)
}

// ListUserTasks godoc
// @Summary      List a user's tasks
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {array}   common.Task
// @Failure      500  {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{id}/tasks [get]
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

// RefreshToken godoc
// @Summary      Refresh tokens
// @Description  Issues a new access token and refresh token. The refresh token must be sent in the Authorization header.
// @Tags         auth
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      201  {object}  common.AuthResponse
// @Failure      500  {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{id}/refresh_token [get]
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

// AddTask godoc
// @Summary      Create a task for a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "User ID"
// @Param        task  body      common.Task  true  "Task data"
// @Success      201   {object}  common.Task
// @Failure      500   {string}  string  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{id}/tasks [post]
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
