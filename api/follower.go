package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lvyonghuan/Ubik-Util/umessenger"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

// This file contains the API handlers for follower call leader.

// followerInitHandel initializes a follower with the given UUID and address.
func followerInitHandel(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
		return
	}

	UUID := c.Query("UUID")
	if UUID == "" {
		engine.Log.Error(errors.New("UUID is empty"))
		errorResponse(c, 400, "UUID is empty")
		return
	}

	addr := c.Query("Addr")
	if addr == "" {
		engine.Log.Error(errors.New("addr is empty"))
		errorResponse(c, 400, "Addr is empty")
		return
	}

	err = engine.AddFollower(addr, UUID)
	if err != nil {
		engine.Log.Error(err)
		errorResponse(c, 400, "add follower err: "+err.Error())
		return
	}

	engine.Log.Debug("Follower init successfully, UUID is: " + UUID)

	successResponse(c, "")
}

func followerPluginListHandel(c *gin.Context) {
	// Get the engine from the context
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
		return
	}

	// Get the UUID from the context
	UUIDAny, isExist := c.Get("UUID")
	if !isExist {
		engine.Log.Error(errors.New("UUID is empty"))
		errorResponse(c, 400, "UUID is empty")
		return
	}
	UUID, ok := UUIDAny.(string)
	if !ok {
		engine.Log.Error(errors.New("UUID is not a string"))
		errorResponse(c, 400, "UUID is not a string")
		return
	}

	// Get the plugin list for the given UUID
	plugins := make(map[string]uplugin.Plugin)
	err = c.BindJSON(&plugins)
	if err != nil {
		engine.Log.Error(errors.New("bind json err: " + err.Error()))
		errorResponse(c, 400, "bind json err: "+err.Error())
		return
	}

	err = engine.AddPlugins(UUID, plugins)
	if err != nil {
		engine.Log.Error(err)
		errorResponse(c, 400, "add plugin err: "+err.Error())
		return
	}

	engine.Log.Debug("Follower add plugin successfully, UUID is: " + UUID)
	successResponse(c, "add plugin success")
}

// TODO
func followerLogHandel(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
		return
	}

	UUIDAny, isExist := c.Get("UUID")
	if !isExist {
		engine.Log.Error(errors.New("UUID is empty"))
		errorResponse(c, 400, "UUID is empty")
		return
	}
	UUID, ok := UUIDAny.(string)
	if !ok {
		engine.Log.Error(errors.New("UUID is not a string"))
		errorResponse(c, 400, "UUID is not a string")
		return
	}

	var envelope umessenger.UEnvelope
	if err := c.BindJSON(&envelope); err != nil {
		engine.Log.Error(errors.New("bind json err: " + err.Error()))
		errorResponse(c, 400, "bind json err: "+err.Error())
		return
	}

	// Print follower's log
	printFollowerLog(UUID, envelope, engine)

	successResponse(c, "add log success")
}
