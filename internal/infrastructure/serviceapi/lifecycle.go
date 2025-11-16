package serviceapi

import (
	"avitostazhko/api"
	"avitostazhko/internal/infrastructure/serviceapi/handlers"
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func RunServer(lc fx.Lifecycle, server *handlers.ReviewService, router *echo.Echo, config *Config) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			api.RegisterHandlers(router, server)
			go func() {
				if err := router.Start(config.Address); err != nil {
					log.Fatalf("server failed to start or finished with error: %v\n", err.Error())
				}
			}()
			log.Printf("Review Server started successfully")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server...")
			return router.Shutdown(ctx)
		},
	})
}
