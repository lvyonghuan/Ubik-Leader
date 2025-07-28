package api

import (
	"Ubik-Leader/engine"

	"github.com/lvyonghuan/Ubik-Util/umessenger"
)

func printFollowerLog(uuid string, envelope umessenger.UEnvelope, e *engine.Engine) {
	message := string(envelope.Message)
	// Print the log message
	e.Log.RecordFollowerLog(uuid, message, envelope.Category)
}
