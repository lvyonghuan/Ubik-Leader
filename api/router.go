package api

import (
	"Ubik-Leader/engine"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lvyonghuan/Ubik-Util/uerr"
)

const (
	engineKey = "engine" //key for the engine in the gin context
)

func InitAPI(engine *engine.Engine) {
	r := gin.Default()
	//set the engine in the context
	r.Use(func(c *gin.Context) {
		c.Set(engineKey, engine)
		c.Next()
	})

	follower := r.Group("/follower")
	{
		follower.GET("/init", followerInitHandel)                                //Init the follower. Then will start heartbeat.
		follower.POST("/list", resetFollowerHeartbeat, followerPluginListHandel) //Post the plugin list

		follower.PUT("/log", resetFollowerHeartbeat, followerLogHandel) //Send the log to the leader
		follower.PUT("/result")                                         //Send the result to the leader
	}

	err := r.Run(":" + engine.Config.Port)
	if err != nil {
		engine.Log.Fatal(err)
		os.Exit(1)
	}
}

// retrieves the engine from the context
func getEngine(c *gin.Context) (*engine.Engine, error) {
	engineVal, isExist := c.Get(engineKey)
	if !isExist {
		return nil, uerr.NewError(errors.New("engine not exist in the context"))
	}
	return engineVal.(*engine.Engine), nil
}

// Reset follower's heartbeat
// If a follower accesses the leader's API, it can be seen as a heartbeat
func resetFollowerHeartbeat(c *gin.Context) {
	e, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
		return
	}

	UUID := c.Query("UUID")
	if UUID == "" {
		e.Log.Error(errors.New("UUID is empty"))
		errorResponse(c, 400, "UUID is empty")
		c.Abort()
		return
	}
	c.Set("UUID", UUID)

	err = e.ResetHeartbeat(UUID)
	if err != nil {
		e.Log.Error(err)
		c.Next()
		return
	}

	c.Next()
}
