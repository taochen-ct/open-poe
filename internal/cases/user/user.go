package user

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewService,
	NewJwtService,
	NewHandler,
	NewRepository,
	NewJwtRepository,
)
