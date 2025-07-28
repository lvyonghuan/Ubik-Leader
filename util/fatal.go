package util

import (
	"github.com/lvyonghuan/Ubik-Util/ulog"
)

// InitFatalErrorHandel
// When error occurs but log is not initialized,
// this function will initialize the log and print the error.
func InitFatalErrorHandel(err error) {
	log := ulog.NewLogWithoutPost(ulog.Debug, true, "./logs")
	log.InitLog()
	log.Fatal(err)
}
