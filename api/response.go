package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lvyonghuan/Ubik-Util/ulog"
)

func response(status int, info any) gin.H {
	return gin.H{
		"status": status,
		"info":   info,
	}
}

func successResponse(c *gin.Context, info any) {
	c.JSON(http.StatusOK, response(200, info))
}

func errorResponse(c *gin.Context, status int, info any) {
	c.JSON(http.StatusOK, response(status, info))
}

func fatalErrHandel(c *gin.Context, err error) {
	l := ulog.NewLogWithoutPost(ulog.Debug, true, "./logs")
	l.Fatal(err)
	c.JSON(http.StatusOK, response(500, "Internal Server Error: "+err.Error()))
	os.Exit(1)
}
