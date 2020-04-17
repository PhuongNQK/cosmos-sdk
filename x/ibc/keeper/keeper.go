package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	wasm "github.com/cosmos/cosmos-sdk/x/ibc/99-wasm/keeper"
)

// Keeper defines each ICS keeper for IBC
type Keeper struct {
	ClientKeeper     client.Keeper
	ConnectionKeeper connection.Keeper
	ChannelKeeper    channel.Keeper
	PortKeeper       port.Keeper
	Router           *port.Router
	WasmKeeper       wasm.Keeper
}

// NewKeeper creates a new ibc Keeper
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, stakingKeeper client.StakingKeeper, scopedKeeper capability.ScopedKeeper,
) *Keeper {
	wasmKeeper := wasm.NewKeeper(cdc, key)
	clientKeeper := client.NewKeeper(cdc, key, stakingKeeper, wasmKeeper)
	connectionKeeper := connection.NewKeeper(cdc, key, clientKeeper)
	portKeeper := port.NewKeeper(scopedKeeper)
	channelKeeper := channel.NewKeeper(cdc, key, clientKeeper, connectionKeeper, portKeeper, scopedKeeper)

	return &Keeper{
		ClientKeeper:     clientKeeper,
		ConnectionKeeper: connectionKeeper,
		ChannelKeeper:    channelKeeper,
		PortKeeper:       portKeeper,
		WasmKeeper:       wasmKeeper,
	}
}

// Set the Router in IBC Keeper and seal it
func (k *Keeper) SetRouter(rtr *port.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
}
