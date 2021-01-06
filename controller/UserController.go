package controller

import (
	"context"
	"davidws/chat"
	"davidws/config"
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

// checkPassword checks equality between password and hash
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

// generatePassword generates a password for user
func generatePassword(password string) string {
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
		expNum:      regexp.MustCompile(`[\d]`),
		expAlpha:    regexp.MustCompile(`[A-z]`),
		expSpecial:  regexp.MustCompile(`[#@!$%^&*()_+|~=\x60{}\[\]:";'<>?,.\/]`),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// checkLogin checks username and password for login
func (uc *UserController) checkLogin(ctx context.Context, username, password string) (userID int, err error) {
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

// checkSignup checks signup process
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

	if len(uc.expSpecial.FindAllStringIndex(username, -1)) > 0 {
		err = errors.New("Username must not include special characters")
		msg = "Username must not include special characters"
		return
	}

	if len(uc.expSpecial.FindAllStringIndex(password, -1)) < 3 || len(uc.expAlpha.FindAllStringIndex(password, -1)) < 3 || len(uc.expNum.FindAllStringIndex(password, -1)) < 2 {
		err = errors.New("Password must have at least 3 alphabetical characters, 2 numbers, and 3 special characters")
		msg = "Password must have at least 3 alphabetical characters, 2 numbers, and 3 special characters"
		return
	}

	user.Hash = generatePassword(password)

	if user.Hash == "" {
		err = errors.New("Hashing process failed")
		msg = "Unknown error"
		return
	}

	return
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
		sendMsgHandler(w, r, "Invalid request", http.StatusBadRequest)
		return
	}

	username, _ := form["username"]
	password, _ := form["password"]

	userID, err := uc.checkLogin(r.Context(), username, password)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		sendMsgHandler(w, r, "Invalid Username or Password", http.StatusUnauthorized)
		return
	}

	sid, err := uc.restartSession(r.Context(), userID, username)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		sendMsgHandler(w, r, "Unknown error", http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{Name: config.SessionIDName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: true, MaxAge: 0})
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

// registerClient registers a new user into the chat
func (uc *UserController) registerClient(w http.ResponseWriter, r *http.Request) (client *chat.Client) {
	sessionInfo := r.Context().Value(ctxtypes.SessionInfo).(model.SessionInfo)

	if sessionInfo.UserID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "User not allowed")
		return
	}

	ws, err := uc.upgrader.Upgrade(w, r, nil)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		sendMsgHandler(w, r, "Unknown error", http.StatusInternalServerError)
		return
	}

	client = &chat.Client{Username: sessionInfo.Username, SessionID: sessionInfo.GetID(), WS: ws, Send: make(chan model.Message), Hub: uc.hub}

	return
}

// restartSession regenerates the sessionID registers it in session storage
func (uc *UserController) restartSession(ctx context.Context, userID int, username string) (sid string, err error) {
	sessionInfo := ctx.Value(ctxtypes.SessionInfo).(model.SessionInfo)

	sessionInfo.Username = username
	sessionInfo.UserID = userID

	sid = sessionInfo.GenerateID(88)

	_, err = uc.sessionRepo.Set(ctx, url.QueryEscape(sid), sessionInfo.String(), time.Duration(time.Minute*15))

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

// Signup manages user sign up process
func (uc *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	form := map[string]string{
		"username": "",
		"password": "",
	}

	err := json.NewDecoder(r.Body).Decode(&form)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		sendMsgHandler(w, r, "Invalid request", http.StatusBadRequest)
		return
	}

	username, _ := form["username"]
	password, _ := form["password"]

	user, msg, err := uc.checkSignup(r.Context(), username, password)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		sendMsgHandler(w, r, msg, http.StatusUnauthorized)
		return
	}

	err = uc.userRepo.Insert(r.Context(), user)

	if err != nil {
		uc.hub.Logger.LogChan <- err.Error()
		sendMsgHandler(w, r, "Unknown error", http.StatusInternalServerError)
		return
	}

	sendMsgHandler(w, r, "Success", 200)
}
