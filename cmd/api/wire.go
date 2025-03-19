//go:build wireinject
// +build wireinject

package main

import (
	"awesomeProject/config"
	"awesomeProject/internal/command"
	commandHandler "awesomeProject/internal/command/handler"
	"awesomeProject/internal/compo"
	"awesomeProject/internal/middleware"
	"awesomeProject/routes"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

// wireApp dependency inject
func wireApp(*config.Configuration, *lumberjack.Logger, *zap.Logger) (*App, func(), error) {
	panic(
		wire.Build(
			compo.ProviderSet,
			middleware.ProviderSet,
			routes.ProviderSet,
			newHttpServer,
			newApp,
		),
	)
}

// wireCommand init application.
func wireCommand(*config.Configuration, *lumberjack.Logger, *zap.Logger) (*command.Command, func(), error) {
	panic(wire.Build(commandHandler.ProviderSet, command.NewCommand))
}
