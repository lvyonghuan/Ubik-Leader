package engine

import "time"

//Handel the action of the workflow, like start, stop and so on.

// RunningWorkflow starts the workflow by preparing all followers and running them.
// It checks the validity of the workflow graph before starting.
// If graph not valid, it returns the node id which is not valid and an error.
func (engine *Engine) RunningWorkflow() (int, error) {
	// Check if the graph is valid before starting the workflow
	// Legal definition: Two start nodes cannot appear
	// in a loop of the workflow graph unless marked as a special start node.
	id, err := engine.graph.CheckGraphValid()
	if err != nil {
		return id, err
	}

	// Prepare the workflow by initializing the start node
	engine.Log.Info("Starting preparation for workflow running")
	// Use a pipe slice to block and wait
	// for all followers to be ready to complete.
	var readyChans []chan struct{}
	// Use an error channel to collect errors from followers
	var errChan = make(chan error, 1)
	// Iterate through all followers and prepare them
	for _, f := range engine.follower.followers {
		readyChans = append(readyChans, f.PreparingFollower(errChan))
	}

	// Waiting for all followers to be ready
	for _, readyChan := range readyChans {
		select {
		case readyChan <- struct{}{}:
			continue
		case err := <-errChan: // If it has error, will return directly
			// FIXME: 还原机制——还原状态。需要从节点补完停止逻辑。
			return 0, err
		}
	}

	time.Sleep(5 * time.Second)

	// Running all followers
	// Use a channel to block and wait for all followers running
	var runningChans []chan struct{}
	engine.Log.Info("All followers are ready, starting workflow running")
	for _, f := range engine.follower.followers {
		runningChans = append(runningChans, f.RunningFollower(errChan))
	}
	// Waiting for all followers to be running
	for _, runningChan := range runningChans {
		select {
		case runningChan <- struct{}{}:
			continue
		case err := <-errChan: // If it has error, will return directly
			//FIXME: 还原机制——还原状态。需要从节点补完停止逻辑。
			return 0, err
		}
	}

	engine.Log.Info("All followers are running.")
	return 0, nil
}
