package controller

import (
	"context"
	"davidws/chat"
	"davidws/ctxtypes"
	"davidws/model"
	"davidws/repo"
	"davidws/security"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

// UserController comprises all user-related endpoints
type UserController struct {
	userRepo    repo.UserRepo
	sessionRepo repo.SessionRepo
	hub         *chat.Hub
	expNum      *regexp.Regexp
	expAlpha    *regexp.Regexp
	expSpecial  *regexp.Regexp
	upgrader    *websocket.Upgrader
}

// NewUserController creates a *UserController
func NewUserController(userRepo repo.UserRepo, sessionRepo repo.SessionRepo, hub *chat.Hub) *UserController {
	return &UserController{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		hub:         hub,
		expNum:      regexp.MustCompile(`[\d]{2,}`),
		expAlpha:    regexp.MustCompile(`[A-z]{3,}`),
		expSpecial:  regexp.MustCompile(`[@!$%^&*()_+|~=\x60{}\[\]:";'<>?,.\/]{3,}`),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Login controls user login operations
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	form := map[string]string{
		"username": "",
		"password": "",
	}

	err := json.NewDecoder(r.Body).Decode(&form)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, "Invalid request", http.StatusBadRequest)
		return
	}

	username, _ := form["username"]
	password, _ := form["password"]

	userID, err := uc.checkUserInfo(r.Context(), username, password)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, "Invalid Username or Password", http.StatusUnauthorized)
		return
	}

	sid, err := uc.restartSession(r.Context(), userID)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, "Unknown error", http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{Name: "sid", Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: true, MaxAge: 0})
}

// Logout erases sessionID from DB
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionInfo := ctx.Value(ctxtypes.SessionInfo).(model.SessionInfo)

	_, err := uc.sessionRepo.Delete(ctx, sessionInfo.GetID())

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Signup manages user signup
func (uc *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	form := map[string]string{
		"username": "",
		"password": "",
	}

	err := json.NewDecoder(r.Body).Decode(&form)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, "Invalid request", http.StatusBadRequest)
		return
	}

	username, _ := form["username"]
	password, _ := form["password"]

	user, msg, err := uc.checkSignup(r.Context(), username, password)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, msg, http.StatusUnauthorized)
		return
	}

	user.Username = username
	user.Hash = uc.generatePassword(password)

	if user.Hash == "" {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, "Unknown error", http.StatusInternalServerError)
		return
	}

	err = uc.userRepo.Insert(r.Context(), user)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, "Unknown error", http.StatusInternalServerError)
		return
	}

	errorHandler(w, r, "Success", 200)
}

// checkSignup checks signup prcss
func (uc *UserController) checkSignup(ctx context.Context, username, password string) (user model.User, msg string, err error) {
	user, err = uc.userRepo.GetByUsername(ctx, username)

	if err == nil {
		err = errors.New("Username already exists")
		msg = "Username already exists"
		return
	}

	err = nil
	user.Username = username

	if len(password) < 8 || len(username) < 5 {
		err = errors.New("Password must be at least 8 characters long and Username at least 5")
		msg = "Password must be at least 8 characters long and Username at least 5"
		return
	}

	if uc.expSpecial.MatchString(username) {
		err = errors.New("Username must not include special characters")
		msg = "Username must not include special characters"
		return
	}

	if !uc.expAlpha.MatchString(password) || !uc.expNum.MatchString(password) || !uc.expSpecial.MatchString(password) {
		err = errors.New("Password must have at least 3 alphabetical characters, 2 numbers, and 3 special characters")
		msg = "Password must have at least 3 alphabetical characters, 2 numbers, and 3 special characters"
		return
	}

	return
}

func (uc *UserController) checkUserInfo(ctx context.Context, username, password string) (userID int, err error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)

	if err != nil {

		return
	}

	same := checkPassword(password, user.Hash)

	if !same {
		err = errors.New("Hashes don't match")
	} else {
		userID = user.ID
	}

	return
}

func (uc *UserController) generatePassword(password string) string {
	done := make(chan string)

	go func(c chan<- string) {
		hashbytes, err := security.Hash(password)

		if err != nil {
			c <- ""
			return
		}
		c <- string(hashbytes)
	}(done)

	return <-done
}

func (uc *UserController) restartSession(ctx context.Context, userID int) (sid string, err error) {
	sessionInfo := ctx.Value(ctxtypes.SessionInfo).(model.SessionInfo)

	sessionInfo.UserID = userID
	sid = sessionInfo.GenerateID(88)

	_, err = uc.sessionRepo.Set(ctx, url.QueryEscape(sid), sessionInfo.UserID, time.Duration(time.Minute*15))

	return
}

func checkPassword(password, hash string) bool {
	done := make(chan bool)

	go func(c chan<- bool) {
		err := security.VerifyPassword(hash, password)

		if err != nil {
			c <- false
			return
		}

		c <- true
	}(done)

	return <-done
}

// registerClient registers a new user into the chat
func (uc *UserController) registerClient(w http.ResponseWriter, r *http.Request) (client *chat.Client) {
	sessionInfo := r.Context().Value(ctxtypes.SessionInfo).(model.SessionInfo)

	if sessionInfo.UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "User not allowed")
		return
	}

	user, err := uc.userRepo.GetByID(r.Context(), sessionInfo.UserID)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, "Unknown error", http.StatusInternalServerError)
		return
	}

	sid := sessionInfo.GetID()

	ws, err := uc.upgrader.Upgrade(w, r, nil)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		errorHandler(w, r, "Unknown error", http.StatusInternalServerError)
		return
	}

	client = &chat.Client{Username: user.Username, SessionID: sid, WS: ws, Send: make(chan model.Message), Hub: uc.hub}

	return
}

// SendMessages contains the main loop where messages are sent
func (uc *UserController) SendMessages(w http.ResponseWriter, r *http.Request) {
	client := uc.registerClient(w, r)

	if client != nil {
		uc.hub.Register <- client
	} else {
		return
	}

	go client.Read()
	go client.Write()
}
