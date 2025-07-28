package engine

import "Ubik-Leader/graph"

// NewRuntimeNode creates a new runtime node in graph
// Will return the id of the new node, and an error if any
func (engine *Engine) NewRuntimeNode(uuid, pluginName, nodeName string) (int, error) {
	return engine.graph.NewRuntimeNode(uuid, pluginName, nodeName)
}

// DeleteRuntimeNode deletes a runtime node from graph
func (engine *Engine) DeleteRuntimeNode(id int) error {
	return engine.graph.DeleteRuntimeNode(id)
}

// UpdateEdge updates an edge in graph
func (engine *Engine) UpdateEdge(producerID, consumerID int, producerPortName, consumerPortName string) error {
	return engine.graph.UpdateEdge(graph.Edge{
		ProducerID:       producerID,
		ConsumerID:       consumerID,
		ProducerPortName: producerPortName,
		ConsumerPortName: consumerPortName,
	})
}

// DeleteEdge deletes an edge in graph
func (engine *Engine) DeleteEdge(producerID, consumerID int, producerPortName, consumerPortName string) error {
	return engine.graph.DeleteEdge(graph.Edge{
		ProducerID:       producerID,
		ConsumerID:       consumerID,
		ProducerPortName: producerPortName,
		ConsumerPortName: consumerPortName,
	})
}
