package api

import (
	"net/http"
	"oauth/Go/common/logger"
	"oauth/Go/common/util"
	"oauth/Go/model"
	db "oauth/migrations/sqlc"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request model.LoginRequest

	if err := util.ReadJsonAndValidate(w, r, &request); err != nil {
		logger.Error("Failed to read/validate request: %v", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

    authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "Auth Header missing or invalid (!Bearer )")
        return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	isValid, err := h.verifyOTPCode(request.VerificationID, request.OTP, token)
    if err != nil {
        logger.Error("Failed to verify OTP: %v", err)
        util.WriteErrorResponse(w, http.StatusInternalServerError, "OTP verification failed")
        return
    }
    if !isValid {
        util.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid OTP")
        return
    }

    user, err := h.db.FindOrCreateUser(r.Context(), db.FindOrCreateUserParams{
        Phone:    pgtype.Text{String: request.Phone, Valid: true},
        Provider: db.ProviderPHONE,
    })
    if err != nil {
        logger.Error("Failed to find/create user: %v", err)
        util.WriteErrorResponse(w, http.StatusInternalServerError, "User creation failed")
        return
    }

	util.WriteSuccessResponse(w, user)
}