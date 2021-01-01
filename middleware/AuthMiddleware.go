package middleware

import (
	"context"
	"davidws/config"
	"davidws/ctxtypes"
	"davidws/model"
	"davidws/repo"
	"davidws/utils"
	"errors"
	"net/http"
	"net/url"
)

// AuthHandler is a handler for sessions
type AuthHandler struct {
	repo   repo.SessionRepo
	logger *utils.Logger
}

// NewAuthHandler creates a *AuthHandler
func NewAuthHandler(repo repo.SessionRepo, logger *utils.Logger) *AuthHandler {
	return &AuthHandler{
		repo:   repo,
		logger: logger,
	}
}

// AuthMiddleware authentication middleware
func (ah AuthHandler) AuthMiddleware(next http.HandlerFunc, auth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var si model.SessionInfo
		c, err := r.Cookie(config.SessionIDName)

		if err != nil || c.Value == "" {
			si.GenerateID(int64(config.SessionIDLength))

			cookie := http.Cookie{Name: config.SessionIDName, Value: url.QueryEscape(si.GetID()), Path: "/", HttpOnly: true, Secure: true, MaxAge: 0}
			http.SetCookie(w, &cookie)

		} else {
			si.SetID(c.Value)
		}

		userID, err := ah.repo.Get(r.Context(), si.GetID())

		err = errors.New("AAAAAA")

		if auth {
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			si.UserID = userID.(string)

		} else {
			if err == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		ctx := context.WithValue(r.Context(), ctxtypes.SessionInfo, si)

		next(w, r.WithContext(ctx))
	}
}
