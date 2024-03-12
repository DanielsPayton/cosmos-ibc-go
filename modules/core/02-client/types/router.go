package types

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// Router is a map from a clientType to a LightClientModule instance.
// The router has a reference to the client store provider (02-client keeper)
// and will register the store provider on a client module upon route registration.
type Router struct {
	routes        map[string]exported.LightClientModule
	storeProvider exported.ClientStoreProvider
}

// NewRouter returns an instance of the Router.
func NewRouter(key storetypes.StoreKey) *Router {
	return &Router{
		routes:        make(map[string]exported.LightClientModule),
		storeProvider: NewStoreProvider(key),
	}
}

// AddRoute adds LightClientModule for a given module name. It returns the Router
// so AddRoute calls can be linked. It will panic if the Router is sealed.
// The store provider will be registered on the light client module.
func (rtr *Router) AddRoute(clientType string, module exported.LightClientModule) *Router {
	if rtr.HasRoute(clientType) {
		panic(fmt.Errorf("route %s has already been registered", module))
	}

	rtr.routes[clientType] = module

	module.RegisterStoreProvider(rtr.storeProvider)
	return rtr
}

// HasRoute returns true if the Router has a module registered or false otherwise.
func (rtr *Router) HasRoute(clientType string) bool {
	_, ok := rtr.routes[clientType]
	return ok
}

// GetRoute returns the LightClientModule registered for the client type
// associated with the clientID.
func (rtr *Router) GetRoute(clientID string) (exported.LightClientModule, bool) {
	clientType, _, err := ParseClientIdentifier(clientID)
	if err != nil {
		return nil, false
	}

	if !rtr.HasRoute(clientType) {
		return nil, false
	}
	return rtr.routes[clientType], true
}
