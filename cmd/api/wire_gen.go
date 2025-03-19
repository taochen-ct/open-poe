// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"awesomeProject/config"
	"awesomeProject/internal/command"
	"awesomeProject/internal/command/handler"
	"awesomeProject/internal/compo"
	"awesomeProject/internal/middleware"
	"awesomeProject/routes"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Injectors from wire.go:

// wireApp dependency inject
func wireApp(configuration *config.Configuration, lumberjackLogger *lumberjack.Logger, zapLogger *zap.Logger) (*App, func(), error) {
	recovery := middleware.NewRecoveryMiddleware(lumberjackLogger)
	cors := middleware.NewCorsMiddleware()
	limiterManager := compo.NewLimiterManager()
	limiter := middleware.NewLimiterMiddleware(limiterManager)
	engine := routes.CreateBaseRouter(recovery, cors, limiter)
	server := newHttpServer(configuration, engine)
	app := newApp(configuration, zapLogger, server)
	return app, func() {
	}, nil
}

// wireCommand init application.
func wireCommand(configuration *config.Configuration, lumberjackLogger *lumberjack.Logger, zapLogger *zap.Logger) (*command.Command, func(), error) {
	exampleHandler := handler.NewExampleHandler(zapLogger)
	commandCommand := command.NewCommand(exampleHandler)
	return commandCommand, func() {
	}, nil
}
