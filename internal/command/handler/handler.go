package handler

import (
	"github.com/google/wire"
)

// ProviderSet command handler providerSet
var ProviderSet = wire.NewSet(NewExampleHandler)
