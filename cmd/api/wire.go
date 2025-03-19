//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"open-poe/config"
	"open-poe/internal/cases/user"
	"open-poe/internal/command"
	commandHandler "open-poe/internal/command/handler"
	"open-poe/internal/compo"
	"open-poe/internal/middleware"
	"open-poe/routes"
)

// wireApp dependency inject
func wireApp(*config.Configuration, *lumberjack.Logger, *zap.Logger) (*App, func(), error) {
	panic(
		wire.Build(
			compo.ProviderSet,
			middleware.ProviderSet,
			user.ProviderSet,
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
