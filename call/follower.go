package call

import (
	"errors"
	"strconv"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

// FIXME: 这里肯定是有问题的，权宜之计。协议应该由从节点指定。
const protocolPrefix = "http://"

type Follower struct {
	addr           string
	uuid           string
	caller         *Caller
	resetHeartbeat chan struct{} // Channel to reset heartbeat
}

// RegisterFollower Register a follower to caller
func (c *Caller) RegisterFollower(uuid, addr string, resetHeartbeat chan struct{}) {
	// Check if the UUID is existing in the caller
	if _, exists := c.Followers[uuid]; exists {
		c.log.Warn("Follower with UUID " + uuid + " already exists in caller, updating address")
	}

	// Register the follower in the caller
	c.Followers[uuid] = Follower{
		addr:           protocolPrefix + addr, //FIXME:参见常量注释
		uuid:           uuid,
		caller:         c,
		resetHeartbeat: resetHeartbeat,
	}
	c.log.Debug("Registered follower in caller: " + uuid + " at " + addr)
}

// GetFollower Get a follower by UUID
func (c *Caller) GetFollower(uuid string) (*Follower, error) {
	if follower, exists := c.Followers[uuid]; exists {
		return &follower, nil
	}
	return nil, uerr.NewError(errors.New("follower with UUID " + uuid + " does not exist"))
}

// When the leader receives a response from follower, can be seen as a heartbeat
func (f Follower) resetFollowerHeartbeat() {
	f.resetHeartbeat <- struct{}{}
}

// GetPluginList Call follower to get a plugin list
// FIXME: Seems unuseful（插件在初始化时已经传递了list）
func (f Follower) GetPluginList() (map[string]uplugin.Plugin, error) {
	uri := f.addr + getPluginList
	f.caller.log.Debug("Calling follower at " + uri + " to get plugin list")

	var plugins map[string]uplugin.Plugin

	status, err := f.caller.callAndUnmarshal("GET", uri, nil, &plugins)
	if err != nil {
		return nil, err
	}

	f.resetFollowerHeartbeat()
	if status != 200 { //TODO 根据状态码分类错误
		return nil, uerr.NewError(errors.New("follower get plugin list error, status code: " + strconv.Itoa(status)))
	}

	f.caller.log.Debug("Received plugin list from follower " + f.uuid)
	return plugins, nil
}

// AddRuntimeNode Add a runtime node to the follower
func (f Follower) AddRuntimeNode(pluginName, nodeName string, id int) error {
	f.caller.log.Debug("Adding runtime node " + strconv.Itoa(id) + " for follower " + f.uuid)

	uri := f.addr + addRuntimeNode + "?id=" + strconv.Itoa(id) + "&plugin_name=" + pluginName + "&node_name=" + nodeName

	var response string
	status, err := f.caller.callAndUnmarshal("PUT", uri, nil, &response)
	if err != nil {
		return err
	}

	f.resetFollowerHeartbeat()
	if status != 200 {
		return uerr.NewError(errors.New("follower add runtime node error, status code: " + strconv.Itoa(status) + ", response: " + response))
	}

	f.caller.log.Debug("Runtime node " + strconv.Itoa(id) + " added successfully for follower " + f.uuid)
	return nil
}

// DeleteRuntimeNode Delete a runtime node from the follower
func (f Follower) DeleteRuntimeNode(id int) error {
	f.caller.log.Debug("Deleting runtime node " + strconv.Itoa(id) + " for follower " + f.uuid)

	uri := f.addr + deleteRuntimeNode + "?id=" + strconv.Itoa(id)

	var response string
	status, err := f.caller.callAndUnmarshal("DELETE", uri, nil, &response)
	if err != nil {
		return err
	}

	f.resetFollowerHeartbeat()
	if status != 200 {
		return uerr.NewError(errors.New("follower delete runtime node error, status code: " + strconv.Itoa(status) + ", response: " + response))
	}

	f.caller.log.Debug("Runtime node " + strconv.Itoa(id) + " deleted successfully for follower " + f.uuid)
	return nil
}

func (f Follower) UpdateEdge(producerID, consumerID int, producerPortName, consumerPortName string) error {
	f.caller.log.Debug("Updating edge from producer " + strconv.Itoa(producerID) + " to consumer " + strconv.Itoa(consumerID) + " for follower " + f.uuid)

	uri := f.addr + updateEdge + "?producer_id=" + strconv.Itoa(producerID) +
		"&consumer_id=" + strconv.Itoa(consumerID) +
		"&producer_port_name=" + producerPortName +
		"&consumer_port_name=" + consumerPortName

	var response string
	status, err := f.caller.callAndUnmarshal("PUT", uri, nil, &response)
	if err != nil {
		return err
	}

	f.resetFollowerHeartbeat()
	if status != 200 {
		return uerr.NewError(errors.New("follower update edge error, status code: " + strconv.Itoa(status) + ", response: " + response))
	}

	f.caller.log.Debug("Edge updated successfully from producer " + strconv.Itoa(producerID) + " to consumer " + strconv.Itoa(consumerID) + " for follower " + f.uuid)
	return nil
}

func (f Follower) DeleteEdge(producerID, consumerID int, producerPortName, consumerPortName string) error {
	f.caller.log.Debug("Deleting edge from producer " + strconv.Itoa(producerID) + " to consumer " + strconv.Itoa(consumerID) + " for follower " + f.uuid)

	uri := f.addr + deleteEdge + "?producer_id=" + strconv.Itoa(producerID) +
		"&consumer_id=" + strconv.Itoa(consumerID) +
		"&producer_port_name=" + producerPortName +
		"&consumer_port_name=" + consumerPortName

	var response string
	status, err := f.caller.callAndUnmarshal("DELETE", uri, nil, &response)
	if err != nil {
		return err
	}

	f.resetFollowerHeartbeat()
	if status != 200 {
		return uerr.NewError(errors.New("follower delete edge error, status code: " + strconv.Itoa(status) + ", response: " + response))
	}

	f.caller.log.Debug("Edge deleted successfully from producer " + strconv.Itoa(producerID) + " to consumer " + strconv.Itoa(consumerID) + " for follower " + f.uuid)
	return nil
}

func (f Follower) PutParams(id int, params map[string]interface{}) error {
	f.caller.log.Debug("Updating parameters for runtime node " + strconv.Itoa(id) + " for follower " + f.uuid)

	uri := f.addr + putParams + "?id=" + strconv.Itoa(id)

	var response string
	status, err := f.caller.callAndUnmarshal("PUT", uri, params, &response)
	if err != nil {
		return err
	}

	f.resetFollowerHeartbeat()
	if status != 200 {
		return uerr.NewError(errors.New("follower update parameters error, status code: " + strconv.Itoa(status) + ", response: " + response))
	}

	f.caller.log.Debug("Parameters updated successfully for runtime node " + strconv.Itoa(id) + " for follower " + f.uuid)
	return nil
}
