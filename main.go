package main

import (
	"context"
	"net/http"
	"time"

	"davidws/chat"
	"davidws/config"
	"davidws/controller"
	"davidws/db"
	"davidws/middleware"
	"davidws/repo/crud"
	"davidws/utils"
)

func main() {
	// Load config
	config.Load()

	// Postgres connection
	pgxPool, err := db.PGXConnect(context.Background())
	utils.FailIfErr(err)
	defer pgxPool.Close()

	err = db.CreateTable(context.Background(), pgxPool)
	utils.FailIfErr(err)

	// Redis connection
	redisClient, err := db.RedisConnect(context.Background())
	utils.FailIfErr(err)
	defer redisClient.Close()

	//var redisClient *redis.Client = nil

	// Repositories for CRUD operations
	sr := crud.NewSessionCRUD(redisClient)
	ur := crud.NewUserCRUD(pgxPool)

	// Concurrent logs
	l := &utils.Logger{LogChan: make(chan interface{})}

	// Chat hub
	h := chat.NewHub(sr, l)

	// Session authentications middleware handler
	ah := middleware.NewAuthHandler(sr, l)
	// User routes controller
	uc := controller.NewUserController(ur, sr, h)

	// Set router
	mux := http.NewServeMux()

	// Backend routes
	mux.HandleFunc("/login", middleware.MethodMiddleware(ah.AuthMiddleware(uc.Login, false), http.MethodPost))
	mux.HandleFunc("/logout", middleware.MethodMiddleware(ah.AuthMiddleware(uc.Logout, true), http.MethodPost))
	mux.HandleFunc("/signup", middleware.MethodMiddleware(ah.AuthMiddleware(uc.Signup, false), http.MethodPost))
	mux.HandleFunc("/chat", ah.AuthMiddleware(uc.SendMessages, true))

	// Frontend route
	mux.Handle("/", http.FileServer(http.Dir("../static/")))

	// Server setup
	s := &http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(time.Second * 5),
		WriteTimeout: time.Duration(time.Second * 5),
		IdleTimeout:  time.Duration(time.Second * 2),
		Addr:         ":8081",
	}

	// Log concurrently
	go l.Logs()
	// Send messages concurrently
	go h.Run(context.Background())

	s.ListenAndServe()
}
