package router

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"social-network-dialogs/internal/logger"
)

var Module = fx.Options(
	fx.Provide(
		NewRouter,
	),
	fx.Invoke(
		InitHooks,
	),
)

func InitHooks(lc fx.Lifecycle, server *Router, logger logger.LoggerInterface) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			logger.Info(fmt.Sprintf("Starting HTTP server at %s", server.Config.HttpServer), nil)
			go func() {
				server.Run()
			}()
			return nil
		},
		//OnStop: func(ctx context.Context) error {
		//	return server.Run(ctx)
		//},
	})
}
