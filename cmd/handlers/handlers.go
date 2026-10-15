package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary      Log in and get tokens
// @Description  Public endpoint. Validates the username and password and returns the user data along with an access token and a refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      common.Login  true  "Login credentials"
// @Success      201          {object}  common.AuthResponse
// @Failure      400          {string}  string  "Invalid request body"
// @Failure      403          {string}  string  "Invalid credentials"
// @Failure      500          {string}  string  "Internal server error"
// @Router       /login [post]
func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	request := &common.Login{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		handler.Logger.Error("Unable to decode request body.", zap.Error(err))
		writeResponse(w, http.StatusBadRequest, "Invalid body")
		return
	}

	user, err := handler.svc.GetUserByUsername(request.Username)
	if err != nil {
		handler.Logger.Error("Unable get user.", zap.Error(err))
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	if user == nil {
		handler.Logger.Error("User not found.")
		writeResponse(w, http.StatusForbidden, "invalid credentials")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		handler.Logger.Error("invalid credentials.")
		writeResponse(w, http.StatusForbidden, "invalid credentials")
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
