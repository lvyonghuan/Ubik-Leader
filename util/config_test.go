package util

import "testing"

func TestReadConfig(t *testing.T) {
	config := ReadConfig("../conf", "config")
	if config.Port == "" {
		t.Error("Port is empty")
	}

	t.Log(config)
}
