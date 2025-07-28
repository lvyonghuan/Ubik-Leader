package engine

import (
	"Ubik-Leader/call"
	"Ubik-Leader/graph"
	"Ubik-Leader/util"
	"errors"
	"os"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ulog"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

type Engine struct {
	Config   util.Config
	caller   *call.Caller
	graph    *graph.Graph
	follower follower

	chs channels

	Log ulog.LeaderLog
}

type channels struct {
	followerHeartbeatError chan error
}

func InitEngine(confPath, configName string) *Engine {
	var engine = &Engine{
		Config: util.ReadConfig(confPath, configName),
	}

	//Initialize log
	engine.Log = ulog.NewLeaderLog(engine.Config.Log.Level, engine.Config.Log.IsSave, engine.Config.Log.SavePath)

	//Initialize caller
	engine.caller = call.InitCaller(engine.Log)

	//Initialize graph
	engine.graph = graph.InitGraph(engine.Log, &engine.follower.followers, engine.caller)

	//Initialize channels
	engine.chs = initChannels()

	//Initialize followers
	err := engine.initFollower()
	if err != nil {
		engine.Log.Fatal(err)
		os.Exit(1)
	}

	engine.Log.Info("Engine initialized successfully")

	return engine
}

func initChannels() channels {
	chs := channels{
		followerHeartbeatError: make(chan error, 1),
	}

	return chs
}

// AddPlugins Register plugins into the follower
// Passing the UUID of the follower and a map of plugins
func (engine *Engine) AddPlugins(UUID string, plugins map[string]uplugin.Plugin) error {
	follower, ok := engine.follower.followers[UUID]
	if !ok {
		return uerr.NewError(errors.New("follower with UUID " + UUID + " does not exist"))
	}

	follower.AddPlugins(plugins)
	return nil
}
