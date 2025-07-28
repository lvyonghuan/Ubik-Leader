package follower

import (
	"errors"
	"time"
)

type heartbeat struct {
	follower *Follower
	timeout  time.Duration
}

func (h *heartbeat) start(overtime int) {
	h.timeout = time.Duration(overtime) * time.Second

	for {
		select {
		case <-h.follower.Chs.StopCh:
			//TODO
			return
		case <-time.After(h.timeout):
			h.follower.log.Debug(h.follower.UUID + " heartbeat timeout")
			h.follower.Chs.HeartbeatErrChan <- errors.New("Heartbeat timeout, UUID: " + h.follower.UUID)
			h.follower.Chs.HeartbeatResetChan <- struct{}{}
		case <-h.follower.Chs.HeartbeatResetChan:
			// Reset the heartbeat timeout
			h.timeout = time.Duration(overtime) * time.Second
		}
	}
}
