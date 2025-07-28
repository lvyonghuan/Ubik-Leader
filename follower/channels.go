package follower

type Channels struct {
	//Input channels

	HeartbeatResetChan chan struct{}
	StopCh             chan struct{}

	//Output channels

	HeartbeatErrChan chan error
}
