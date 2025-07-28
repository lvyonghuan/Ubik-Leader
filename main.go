package main

import (
	"Ubik-Leader/api"
	"Ubik-Leader/engine"
)

const configName = "config"

func main() {
	e := engine.InitEngine("./conf/", configName)
	api.InitAPI(e)
}
