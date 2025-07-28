package engine

import (
	f "Ubik-Leader/follower"
	"errors"
	"net"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ujson"
)

type follower struct {
	conn      *net.UDPConn
	followers f.Followers
}
type heartbeatPacket struct {
	UUID string
}

func (engine *Engine) initFollower() error {
	engine.follower.followers = make(f.Followers)
	var err error
	engine.follower.conn, err = newHeartbeatListener(engine.Config.Port)
	if err != nil {
		return uerr.NewError(err)
	}

	engine.Log.Debug("Follower Controller initialized")

	// Start the heartbeat listener
	go engine.startHeartbeatListener(engine.follower.conn)

	return nil
}

// AddFollower Add a follower to the engine
func (engine *Engine) AddFollower(addr, UUID string) error {
	if _, ok := engine.follower.followers[UUID]; ok {
		return uerr.NewError(errors.New("follower UUID " + UUID + " already exists"))
	}

	chs := f.Channels{
		HeartbeatResetChan: make(chan struct{}, 1),
		StopCh:             make(chan struct{}, 1),
		HeartbeatErrChan:   engine.chs.followerHeartbeatError,
	}

	engine.follower.followers[UUID] = f.InitFollower(addr, UUID, chs, engine.Log, engine.caller, engine.Config.Heartbeat.Overtime)

	engine.Log.Debug("Follower " + UUID + " added successfully")
	return nil
}

func newHeartbeatListener(port string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		return nil, uerr.NewError(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, uerr.NewError(err)
	}

	return conn, nil
}

func (engine *Engine) startHeartbeatListener(conn *net.UDPConn) {
	// Start the heartbeat listener
	for {
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			engine.Log.Error(uerr.NewError(errors.New("Error reading from UDP connection: " + err.Error())))
			continue
		}
		// Process the heartbeat message
		go engine.processHeartbeatMessage(buf[:n])
	}
}

func (engine *Engine) processHeartbeatMessage(data []byte) {
	// Unmarshal the heartbeat message
	var packet heartbeatPacket
	err := ujson.Unmarshal(data, &packet)
	if err != nil {
		engine.Log.Error(uerr.NewError(errors.New("Error unmarshalling heartbeat packet: " + err.Error())))
		return
	}

	err = engine.ResetHeartbeat(packet.UUID)
	if err != nil {
		engine.Log.Error(err)
	}
}

func (engine *Engine) ResetHeartbeat(uuid string) error {
	// Process the heartbeat packet
	follower, ok := engine.follower.followers[uuid]
	if !ok {
		// FIXME 在从机向主机注册前的访问（理论上全部为日志打印）不可能找到UUID，但不可能不打印日志。
		// 然而理论上来说，从机通过网络访问主机时，等同于可以完成注册。所以这里会非常头大。
		return nil
	}

	go func() {
		follower.Chs.HeartbeatResetChan <- struct{}{}
	}()

	return nil
}
