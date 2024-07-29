package main

import (
	"log/slog"
	"net/http"
	"router/api"
	"router/config"
	"router/handler"
	"router/middleware"
	"router/mongodb"

	"github.com/gin-gonic/gin"
)

const configPath string = "config.toml"

func main() {
	// Load config
	conf := config.NewConfig(configPath)
	if conf == nil {
		return
	}
	// Connect DB
	dbClient := mongodb.NewDBClient(conf.DB.DBURI, conf.DB.DBName)
	// New Router
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	if conf.Router.EnableStaticFS {
		router.StaticFS("/maplestory/", http.Dir(conf.Router.StaticFilePath))
	}
	// For Auth user
	authGroup := router.Group(api.PathAuthGroup)
	authGroup.Use(middleware.JwtAuthMiddleware(conf.Token.AccessTokenSignKey))
	// Account
	accountHandler := handler.NewAccountHandler(conf, dbClient)
	router.POST(api.PathAccount, accountHandler.Signup)
	router.PATCH(api.PathAccountPassword, accountHandler.Retrieve)
	authGroup.GET(api.PathAuthAccount, accountHandler.GetAccount)
	authGroup.PATCH(api.PathAuthAccountPassword, accountHandler.UpdatePassword)
	authGroup.PATCH(api.PathAuthAccountSecondPassword, accountHandler.UpdateSecondPassword)
	// JWT
	jwtHandler := handler.NewJWTHandler(conf, dbClient)
	router.POST(api.PathJWT, jwtHandler.Login)
	router.PATCH(api.PathJWT, jwtHandler.Renew)
	// Game
	gameHandler := handler.NewGameHandler(conf, dbClient)
	authGroup.GET(api.PathAuthGameSkipSDOAuth, gameHandler.SkipSDOAuth)
	authGroup.GET(api.PathAuthGameKick, gameHandler.KickGameClient)
	// Bind router to http server
	s := &http.Server{
		Addr:    conf.Router.Addr,
		Handler: router,
	}
	err := s.ListenAndServe()
	if err != nil {
		slog.Error("Failed to run router", "err", err)
		return
	}
	slog.Info("HTTP server run at", "addr", s.Addr)
}
