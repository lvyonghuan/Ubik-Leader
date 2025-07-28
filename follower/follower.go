package follower

import (
	"Ubik-Leader/call"

	"github.com/lvyonghuan/Ubik-Util/ulog"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

type Followers map[string]*Follower

type Follower struct {
	UUID    string //Unique identifier for the follower
	Addr    string // Addr of the follower
	Plugins map[string]uplugin.Plugin

	Chs Channels

	heartbeat *heartbeat

	log    ulog.Log
	caller *call.Caller
}

func InitFollower(addr, UUID string, channels Channels, log ulog.Log, caller *call.Caller, overtime int) *Follower {
	f := &Follower{
		Addr:      addr,
		UUID:      UUID,
		log:       log,
		caller:    caller,
		Chs:       channels,
		Plugins:   make(map[string]uplugin.Plugin),
		heartbeat: &heartbeat{},
	}
	f.heartbeat.follower = f

	f.log.Debug("Follower initialized")

	// Initialize the heartbeat
	go f.heartbeat.start(overtime)

	f.caller.RegisterFollower(UUID, addr, f.Chs.HeartbeatResetChan)

	return f
}

// AddPlugins is a method for adding plugins to the follower.
func (f *Follower) AddPlugins(plugins map[string]uplugin.Plugin) {
	f.Plugins = plugins
	f.log.Debug("Plugins added to follower " + f.UUID)
}
