//go:generate protoc --proto_path=./internal/server/proto/ --go_out=./internal/server/proto/ --go-grpc_out=./internal/server/proto/ ./internal/server/proto/dialogs.proto ./internal/server/proto/pagination.proto
package main

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"social-network-dialogs/internal/config"
	"social-network-dialogs/internal/database"
	"social-network-dialogs/internal/dialog"
	"social-network-dialogs/internal/logger"
	"social-network-dialogs/internal/server"
)

func init() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func main() {
	fxContainer := fx.New(
		logger.Module,
		config.Module,
		database.Module,
		dialog.Module,
		server.Module,
	)

	fxContainer.Run()
}
