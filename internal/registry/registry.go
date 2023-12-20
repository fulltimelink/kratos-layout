package registry

import "github.com/google/wire"

var RegistryProviderSet = wire.NewSet(NewConsulRegistry)
