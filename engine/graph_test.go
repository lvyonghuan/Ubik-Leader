package engine_test

import (
	"Ubik-Leader/api"
	"Ubik-Leader/engine"
	"testing"
	"time"
)

const testFollowerUUID = "cfa98be0-7452-49dd-b769-4706daa7228b"

func initTest() *engine.Engine {
	e := engine.InitEngine("../conf/", "config_test")
	go api.InitAPI(e)

	return e
}

func TestNewRuntimeNode(t *testing.T) {
	e := initTest()

	time.Sleep(5 * time.Second)

	// Create a new runtime node
	idA, err := e.NewRuntimeNode(testFollowerUUID, "AddNum", "startNode")
	if err != nil {
		t.Fatalf("Failed to create new runtime node: %v", err)
	} else {
		t.Logf("New runtime node created with ID: %d", idA)
	}

	idB, err := e.NewRuntimeNode(testFollowerUUID, "AddNum", "selfIncreasingNode")
	if err != nil {
		t.Fatalf("Failed to create new runtime node: %v", err)
	} else {
		t.Logf("New runtime node created with ID: %d", idB)
	}

	idC, err := e.NewRuntimeNode(testFollowerUUID, "AddNum", "sumNode")
	if err != nil {
		t.Fatalf("Failed to create new runtime node: %v", err)
	} else {
		t.Logf("New runtime node created with ID: %d", idC)
	}
}

func createRuntimeNodeForTest(uuid string, e *engine.Engine) (int, int, int, int) {
	idA, err := e.NewRuntimeNode(uuid, "AddNum", "startNode")
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	idB, err := e.NewRuntimeNode(uuid, "AddNum", "selfIncreasingNode")
	if err != nil {
		panic(err)
	}

	idC, err := e.NewRuntimeNode(uuid, "AddNum", "sumNode")
	if err != nil {
		panic(err)
	}

	idD, err := e.NewRuntimeNode(uuid, "AddNum", "selfIncreasingNode")

	return idA, idB, idC, idD
}

func TestDeleteRuntimeNode(t *testing.T) {
	e := initTest()
	time.Sleep(5 * time.Second)
	idA, idB, idC, _ := createRuntimeNodeForTest(testFollowerUUID, e)
	// Delete the runtime nodes
	err := e.DeleteRuntimeNode(idA)
	if err != nil {
		t.Fatalf("Failed to delete runtime node A: %v", err)
	} else {
		t.Logf("Runtime node A with ID %d deleted successfully", idA)
	}

	err = e.DeleteRuntimeNode(idB)
	if err != nil {
		t.Fatalf("Failed to delete runtime node B: %v", err)
	} else {
		t.Logf("Runtime node B with ID %d deleted successfully", idB)
	}

	err = e.DeleteRuntimeNode(idC)
	if err != nil {
		t.Fatalf("Failed to delete runtime node C: %v", err)
	} else {
		t.Logf("Runtime node C with ID %d deleted successfully", idC)
	}
}

func TestCreateEdge(t *testing.T) {
	e := initTest()

	time.Sleep(5 * time.Second)

	idA, idB, idC, idD := createRuntimeNodeForTest(testFollowerUUID, e)

	// Create edges between nodes
	err := e.UpdateEdge(idA, idC, "num_a", "num_a")
	if err != nil {
		t.Fatalf("Failed to create edge from A to C: %v", err)
	} else {
		t.Logf("Edge created successfully from node A (ID: %d) to node C (ID: %d)", idA, idC)
	}

	err = e.UpdateEdge(idA, idB, "num_a", "input")
	if err != nil {
		t.Fatalf("Failed to create edge from A to B: %v", err)
	} else {
		t.Logf("Edge created successfully from node A (ID: %d) to node B (ID: %d)", idA, idB)
	}

	err = e.UpdateEdge(idB, idC, "num_b", "num_b")
	if err != nil {
		t.Fatalf("Failed to create edge from B to C: %v", err)
	} else {
		t.Logf("Edge created successfully from node B (ID: %d) to node C (ID: %d)", idB, idC)
	}

	err = e.UpdateEdge(idA, idD, "cycle_num", "input")
	if err != nil {
		t.Fatalf("Failed to create edge from C to D: %v", err)
	} else {
		t.Logf("Edge created successfully from node C (ID: %d) to node D (ID: %d)", idC, idD)
	}

	err = e.UpdateEdge(idD, idA, "num_b", "current_cycle_num")
	if err != nil {
		t.Fatalf("Failed to create edge from D to A: %v", err)
	} else {
		t.Logf("Edge created successfully from node D (ID: %d) to node A (ID: %d)", idD, idA)
	}

	err = e.UpdateEdge(idC, idA, "sum", "num_a")
	if err != nil {
		t.Fatalf("Failed to create edge from C to A: %v", err)
	} else {
		t.Logf("Edge created successfully from node C (ID: %d) to node A (ID: %d)", idC, idA)
	}
}

func createEdgeForTest(e engine.Engine, idA, idB, idC, idD int) {
	// Create edges between nodes
	err := e.UpdateEdge(idA, idC, "num_a", "num_a")
	if err != nil {
		panic(err)
	}

	err = e.UpdateEdge(idA, idB, "num_a", "input")
	if err != nil {
		panic(err)
	}

	err = e.UpdateEdge(idB, idC, "num_b", "num_b")
	if err != nil {
		panic(err)
	}

	err = e.UpdateEdge(idA, idD, "cycle_num", "input")
	if err != nil {
		panic(err)
	}

	err = e.UpdateEdge(idD, idA, "num_b", "current_cycle_num")
	if err != nil {
		panic(err)
	}

	err = e.UpdateEdge(idC, idA, "sum", "num_a")
	if err != nil {
		panic(err)
	}
}

func TestDeleteEdge(t *testing.T) {
	e := initTest()
	time.Sleep(5 * time.Second)

	idA, idB, idC, idD := createRuntimeNodeForTest(testFollowerUUID, e)
	createEdgeForTest(*e, idA, idB, idC, idD)

	// Delete edges between nodes
	err := e.DeleteEdge(idA, idC, "num_a", "num_a")
	if err != nil {
		t.Fatalf("Failed to delete edge from A to C: %v", err)
	} else {
		t.Logf("Edge deleted successfully from node A (ID: %d) to node C (ID: %d)", idA, idC)
	}

	err = e.DeleteEdge(idA, idB, "num_a", "input")
	if err != nil {
		t.Fatalf("Failed to delete edge from A to B: %v", err)
	} else {
		t.Logf("Edge deleted successfully from node A (ID: %d) to node B (ID: %d)", idA, idB)
	}

	err = e.DeleteEdge(idB, idC, "num_b", "num_b")
	if err != nil {
		t.Fatalf("Failed to delete edge from B to C: %v", err)
	} else {
		t.Logf("Edge deleted successfully from node B (ID: %d) to node C (ID: %d)", idB, idC)
	}

	err = e.DeleteEdge(idA, idD, "cycle_num", "input")
	if err != nil {
		t.Fatalf("Failed to delete edge from A to D: %v", err)
	} else {
		t.Logf("Edge deleted successfully from node A (ID: %d) to node D (ID: %d)", idA, idD)
	}

	err = e.DeleteEdge(idD, idA, "num_b", "current_cycle_num")
	if err != nil {
		t.Fatalf("Failed to delete edge from D to A: %v", err)
	} else {
		t.Logf("Edge deleted successfully from node D (ID: %d) to node A (ID: %d)", idD, idA)
	}

	err = e.DeleteEdge(idC, idA, "sum", "num_a")
	if err != nil {
		t.Fatalf("Failed to delete edge from C to A: %v", err)
	} else {
		t.Logf("Edge deleted successfully from node C (ID: %d) to node A (ID: %d)", idC, idA)
	}
}

func TestDeleteRuntimeNodeWithEdge(t *testing.T) {
	e := initTest()
	time.Sleep(5 * time.Second)

	idA, idB, idC, idD := createRuntimeNodeForTest(testFollowerUUID, e)
	createEdgeForTest(*e, idA, idB, idC, idD)

	// Delete runtime nodes with edges
	err := e.DeleteRuntimeNode(idA)
	if err != nil {
		t.Fatalf("Failed to delete runtime node A: %v", err)
	} else {
		t.Logf("Runtime node A with ID %d deleted successfully", idA)
	}

	err = e.DeleteRuntimeNode(idB)
	if err != nil {
		t.Fatalf("Failed to delete runtime node B: %v", err)
	} else {
		t.Logf("Runtime node B with ID %d deleted successfully", idB)
	}

	err = e.DeleteRuntimeNode(idC)
	if err != nil {
		t.Fatalf("Failed to delete runtime node C: %v", err)
	} else {
		t.Logf("Runtime node C with ID %d deleted successfully", idC)
	}

	err = e.DeleteRuntimeNode(idD)
	if err != nil {
		t.Fatalf("Failed to delete runtime node D: %v", err)
	} else {
		t.Logf("Runtime node D with ID %d deleted successfully", idD)
	}
}

func TestPutParams(t *testing.T) {
	e := initTest()
	time.Sleep(5 * time.Second)

	idA, idB, idC, idD := createRuntimeNodeForTest(testFollowerUUID, e)
	createEdgeForTest(*e, idA, idB, idC, idD)

	// Put parameters into runtime nodes
	params := make(map[string]any)
	params["init_num"] = 0
	params["cycle_num"] = 10

	err := e.PutParams(idA, params)
	if err != nil {
		t.Fatalf("Failed to put parameters into runtime node A: %v", err)
	} else {
		t.Logf("Parameters put successfully into runtime node A with ID: %d", idA)
	}
}

func putParamsForTest(e engine.Engine, id int, params map[string]any) {
	err := e.PutParams(id, params)
	if err != nil {
		panic(err)
	}
}
