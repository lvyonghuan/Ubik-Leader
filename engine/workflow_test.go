package engine_test

import (
	"testing"
	"time"
)

func TestRunWorkflow(t *testing.T) {
	e := initTest()
	time.Sleep(5 * time.Second)

	idA, idB, idC, idD := createRuntimeNodeForTest(testFollowerUUID, e)
	createEdgeForTest(*e, idA, idB, idC, idD)

	params := make(map[string]any)
	params["init_num"] = 0
	params["cycle_num"] = 10
	putParamsForTest(*e, idA, params)

	err := e.RunningWorkflow()
	if err != nil {
		t.Errorf("RunningWorkflow failed: %v", err)
	}

	time.Sleep(20 * time.Second)

	t.Logf("Workflow completed successfully")
}
