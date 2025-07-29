package graph

import (
	"Ubik-Leader/call"
	"Ubik-Leader/follower"
	"errors"
	"strconv"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ujson"
	"github.com/lvyonghuan/Ubik-Util/ulog"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

type Graph struct {
	log       ulog.Log
	followers *follower.Followers
	caller    *call.Caller

	mountPlugins map[string]uplugin.Plugin //When a plugin has a runtime node, it will be mounted to the graph.
	runtimeNodes map[int]RuntimeNode

	currentID int
}

type RuntimeNode struct {
	ID       int    `json:"id"`
	NodeName string `json:"node_name"`

	params uplugin.Params

	FollowerUUID string            `json:"follower_uuid"`
	PluginName   string            `json:"plugin_name"`
	OutputEdges  map[string][]Edge `json:"output_edges"`
	inputEdges   map[string][]Edge
	addr         string //Plugin's addr
}

type Edge struct {
	ProducerID       int    `json:"producer_id"`
	ConsumerID       int    `json:"consumer_id"`
	ProducerPortName string `json:"producer_port_name"`
	ConsumerPortName string `json:"consumer_port_name"`
	Addr             string `json:"addr"` //The addr point to the consumer PLUGIN.
}

func InitGraph(log ulog.Log, followers *follower.Followers, caller *call.Caller) *Graph {
	g := &Graph{
		log:          log,
		followers:    followers,
		caller:       caller,
		mountPlugins: make(map[string]uplugin.Plugin),
		runtimeNodes: make(map[int]RuntimeNode),
		currentID:    0,
	}

	g.log.Debug("Graph initialized")

	return g
}

func (g *Graph) NewRuntimeNode(uuid, pluginName, nodeName string) (int, error) {
	//Check if the follower exists
	f, isExist := (*g.followers)[uuid]
	if !isExist {
		return 0, uerr.NewError(errors.New("Follower " + uuid + " not exist"))
	}

	//Check the plugin and node exist
	plugin, isExist := f.Plugins[pluginName]
	if !isExist {
		return 0, uerr.NewError(errors.New("Plugin " + pluginName + " not exist"))
	}
	node, isExist := plugin.Nodes[nodeName]
	if !isExist {
		return 0, uerr.NewError(errors.New("Node " + nodeName + " not exist in plugin" + pluginName))
	}

	//Add the node to the runtime nodes
	g.currentID++
	runtimeNode := RuntimeNode{
		ID:           g.currentID,
		NodeName:     nodeName,
		params:       make(uplugin.Params),
		FollowerUUID: uuid,
		PluginName:   pluginName,
		OutputEdges:  make(map[string][]Edge),
		inputEdges:   make(map[string][]Edge),
		addr:         plugin.Addr,
	}

	g.runtimeNodes[g.currentID] = runtimeNode

	// Init params key
	for key := range node.Params {
		// TODO: 默认初始值
		runtimeNode.params[key] = nil // Initialize with nil
	}

	//check plugin mount status
	if _, isExist := g.mountPlugins[pluginName]; !isExist {
		g.mountPlugins[pluginName] = plugin
	}

	cFollower, err := g.caller.GetFollower(f.UUID)
	if err != nil {
		return 0, err
	}
	//Call the follower to add the runtime node
	err = cFollower.AddRuntimeNode(pluginName, nodeName, g.currentID)
	if err != nil {
		return 0, err
	}

	//return the id
	return g.currentID, nil
}

func (g *Graph) DeleteRuntimeNode(id int) error {
	var rNode RuntimeNode
	var isExist bool

	//Check if the node exists
	if rNode, isExist = g.runtimeNodes[id]; !isExist {
		return uerr.NewError(errors.New("Runtime node " + strconv.Itoa(id) + " not exist"))
	}

	//Delete the output edges
	for _, edges := range rNode.OutputEdges {
		for _, edge := range edges {
			consumerNode, isExist := g.runtimeNodes[edge.ConsumerID]
			if !isExist {
				return uerr.NewError(errors.New("Consumer node " + strconv.Itoa(edge.ConsumerID) + " not exist"))
			}

			err := g.deleteEdge(edge, rNode, consumerNode)
			if err != nil {
				g.log.Error(err)
				continue
			}
		}
	}

	//Delete input edges
	for _, edges := range rNode.inputEdges {
		for _, edge := range edges {
			producerNode, isExist := g.runtimeNodes[edge.ProducerID]
			if !isExist {
				return uerr.NewError(errors.New("Producer node " + strconv.Itoa(edge.ProducerID) + " not exist"))
			}

			err := g.deleteEdge(edge, producerNode, rNode)
			if err != nil {
				g.log.Error(err)
				continue
			}
		}
	}

	//Delete the node from the runtime nodes
	delete(g.runtimeNodes, id)

	cFollower, err := g.caller.GetFollower(rNode.FollowerUUID)
	if err != nil {
		return err
	}
	//Call the follower to delete the runtime node
	err = cFollower.DeleteRuntimeNode(id)

	return nil
}

func (g *Graph) UpdateEdge(edge Edge) error {
	edge.Addr = g.runtimeNodes[edge.ConsumerID].addr // Set the addr to the consumer node's addr

	var producerNode, consumerNode RuntimeNode
	var producerPort, consumerPort uplugin.Port
	var isExist bool

	//Check if the node exists
	if producerNode, isExist = g.runtimeNodes[edge.ProducerID]; !isExist {
		return uerr.NewError(errors.New("Producer node " + strconv.Itoa(edge.ConsumerID) + " not exist"))
	}
	if consumerNode, isExist = g.runtimeNodes[edge.ConsumerID]; !isExist {
		return uerr.NewError(errors.New("Consumer node " + strconv.Itoa(edge.ConsumerID) + " not exist"))
	} else {
		//Get consumer node's addr form it's plugin's metadata
		consumerPlugin, isExist := (*g.followers)[consumerNode.FollowerUUID].Plugins[consumerNode.PluginName]
		if !isExist {
			return uerr.NewError(errors.New("Consumer plugin " + consumerNode.PluginName + " not exist"))
		}

		edge.Addr = consumerPlugin.Addr
	}

	//Check if the port exists
	if producerPort, isExist = (*g.followers)[producerNode.FollowerUUID].Plugins[producerNode.PluginName].Nodes[producerNode.NodeName].Output[edge.ProducerPortName]; !isExist {
		return uerr.NewError(errors.New("Producer port " + edge.ProducerPortName + " not exist"))
	}
	if consumerPort, isExist = (*g.followers)[consumerNode.FollowerUUID].Plugins[consumerNode.PluginName].Nodes[consumerNode.NodeName].Input[edge.ConsumerPortName]; !isExist {
		return uerr.NewError(errors.New("Consumer port " + edge.ConsumerPortName + " not exist"))
	}

	//Check if we can establish a link from the producer port to the consumer port
	if !(producerPort.Attribute == consumerPort.Attribute) {
		return uerr.NewError(errors.New("Producer port " + edge.ProducerPortName + " and consumer port " + edge.ConsumerPortName + " cannot be linked.\n" +
			"Producer port attribute: " + producerPort.Attribute + "\n" +
			"Consumer port attribute: " + consumerPort.Attribute))
	}

	g.runtimeNodes[edge.ProducerID].OutputEdges[edge.ProducerPortName] = append(g.runtimeNodes[edge.ProducerID].OutputEdges[edge.ProducerPortName], edge)
	g.runtimeNodes[edge.ConsumerID].inputEdges[edge.ConsumerPortName] = append(g.runtimeNodes[edge.ConsumerID].inputEdges[edge.ConsumerPortName], edge)

	cFollower, err := g.caller.GetFollower(producerNode.FollowerUUID)
	if err != nil {
		return err
	}
	err = cFollower.UpdateEdge(edge.ProducerID, edge.ConsumerID, edge.ProducerPortName, edge.ConsumerPortName, edge.Addr)

	return nil
}

func (g *Graph) DeleteEdge(edge Edge) error {
	var producerNode, consumerNode RuntimeNode
	var isExist bool
	//Check if the node exists
	if producerNode, isExist = g.runtimeNodes[edge.ProducerID]; !isExist {
		return uerr.NewError(errors.New("Producer node " + strconv.Itoa(edge.ConsumerID) + " not exist"))
	}
	if consumerNode, isExist = g.runtimeNodes[edge.ConsumerID]; !isExist {
		return uerr.NewError(errors.New("Consumer node " + strconv.Itoa(edge.ConsumerID) + " not exist"))
	}

	return g.deleteEdge(edge, producerNode, consumerNode)
}

func (g *Graph) deleteEdge(edge Edge, producerNode, consumerNode RuntimeNode) error {
	var producerEdgeIsExist = false
	//Check if the edge exists
	for i, e := range producerNode.OutputEdges[edge.ProducerPortName] {
		if e.ConsumerID == edge.ConsumerID && e.ConsumerPortName == edge.ConsumerPortName {
			producerNode.OutputEdges[edge.ProducerPortName] = append(producerNode.OutputEdges[edge.ProducerPortName][:i], producerNode.OutputEdges[edge.ProducerPortName][i+1:]...)
			cFollower, err := g.caller.GetFollower(producerNode.FollowerUUID)
			if err != nil {
				return err
			}
			err = cFollower.DeleteEdge(edge.ProducerID, edge.ConsumerID, edge.ProducerPortName, edge.ConsumerPortName)
			if err != nil {
				return err
			}

			producerEdgeIsExist = true
		}
	}

	if !producerEdgeIsExist {
		return uerr.NewError(errors.New("producer edge not exist"))
	}

	for i, e := range consumerNode.inputEdges[edge.ConsumerPortName] {
		if e.ProducerID == edge.ProducerID && e.ProducerPortName == edge.ProducerPortName {
			consumerNode.inputEdges[edge.ConsumerPortName] = append(consumerNode.inputEdges[edge.ConsumerPortName][:i], consumerNode.inputEdges[edge.ConsumerPortName][i+1:]...)
			cFollower, err := g.caller.GetFollower(consumerNode.FollowerUUID)
			if err != nil {
				g.log.Error(err)
				continue
			}
			err = cFollower.DeleteEdge(edge.ProducerID, edge.ConsumerID, edge.ProducerPortName, edge.ConsumerPortName)
			if err != nil {
				g.log.Error(err)
				continue
			}
			return nil
		}
	}

	//If the edge does not exist (not found in the output edges), return an error
	return uerr.NewError(errors.New("consumer edge not exist"))
}

// PutParams updates the parameters of a runtime node.
func (g *Graph) PutParams(id int, params map[string]any) error {
	// Check if the node exists
	rNode, isExist := g.runtimeNodes[id]
	if !isExist {
		return errors.New("Runtime node " + strconv.Itoa(id) + " not exist")
	}

	// Check if the params are valid
	for key, value := range params {
		if _, isExist := rNode.params[key]; !isExist {
			return uerr.NewError(errors.New("Params port " + strconv.Itoa(id) + " not exist"))
		} else {
			// Update the param value
			var err error
			rNode.params[key], err = ujson.Marshal(value)
			if err != nil {
				return err
			}
		}
	}

	// Update the runtime node's params
	// Get Caller for the follower
	cFollower, err := g.caller.GetFollower(rNode.FollowerUUID)
	if err != nil {
		return err
	}

	// Call the follower to update the params
	err = cFollower.PutParams(id, params)
	if err != nil {
		return err
	}

	g.runtimeNodes[id] = rNode //Update the runtime node in the graph
	g.log.Debug("Updated params for runtime node" + strconv.Itoa(id))
	return nil
}

// CheckGraphValid checks if the graph is valid.
// Legal definition: Two start nodes cannot appear
// in a loop of the workflow graph unless marked as a special start node.
// TODO: use DFS to check if the graph is valid
func (g *Graph) CheckGraphValid() error {
	g.log.Debug("Start checking graph validity")

	g.log.Debug("Graph validity check completed.")
	return nil
}
