package call

const (
	apiRoot       = "/api"            //Will splice behind the port
	getPluginList = apiRoot + "/list" //Get the list of plugins

	prepare = apiRoot + "/prepare" //Ask the follower to prepare
	run     = apiRoot + "/run"     //Ask the follower to run. Must after prepare.

	node              = "/node"
	addRuntimeNode    = apiRoot + node + "/"      //Add a new runtime node. Use put.
	deleteRuntimeNode = apiRoot + node + "/"      //Delete a runtime node. Use delete.
	updateEdge        = apiRoot + node + "/edge"  //Update a edge. Use put.
	deleteEdge        = apiRoot + node + "/edge"  //Delete a edge. Use delete.
	putParams         = apiRoot + node + "/param" //Update the parameters of a node. Use put.
)
