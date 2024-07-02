package main

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"social-network-dialogs/internal/config"
	"social-network-dialogs/internal/database"
	"social-network-dialogs/internal/dialog"
	"social-network-dialogs/internal/logger"
	"social-network-dialogs/internal/rest/handler"
	"social-network-dialogs/internal/rest/router"
	"social-network-dialogs/internal/token"
)

func init() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func main() {
	fxContainer := fx.New(
		logger.Module,
		config.Module,
		database.Module,
		handler.Module,
		router.Module,
		token.Module,
		dialog.Module,
	)

	fxContainer.Run()
}
