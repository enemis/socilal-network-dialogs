package router

import (
	"github.com/gin-gonic/gin"
	"social-network-dialogs/internal/config"
	"social-network-dialogs/internal/middleware"
	"social-network-dialogs/internal/rest/handler"
	"social-network-dialogs/internal/token"
)

type Router struct {
	Router *gin.Engine
	Config *config.Config
}

func NewRouter(dialogHandler *handler.DialogHandler, config *config.Config, generator token.PasswordGenerator) *Router {
	router := gin.New()
	router.Use(middleware.AuthRequired(generator))

	router.POST("/dialog/send", dialogHandler.SendMessage)
	router.GET("/dialog/:user_id/list", dialogHandler.GetDialog)

	return &Router{
		Router: router,
		Config: config,
	}
}

func (router *Router) Run() error {
	return router.Router.Run(router.Config.HttpServer)
}

func (router *Router) Stop() error {
	return router.Router.Run()
}
