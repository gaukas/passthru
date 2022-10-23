package config_test

import (
	"testing"

	"github.com/gaukas/passthru/config"
)

var (
	conf *config.Config
)

func TestConfig(t *testing.T) {
	testLoadConfig(t)
	testVersion(t)
	testServers(t)
}

func testLoadConfig(t *testing.T) {
	var err error
	conf, err = config.LoadConfig("./test.json")
	if err != nil {
		t.Errorf("failed to load config: %v", err)
		return
	}
}

func testVersion(t *testing.T) {
	if conf.MinVersion.Major != 3 {
		t.Errorf("incorrect major version: %d", conf.MinVersion.Major)
	}
	if conf.MinVersion.Minor != 2 {
		t.Errorf("incorrect minor version: %d", conf.MinVersion.Minor)
	}
	if conf.MinVersion.Patch != 1 {
		t.Errorf("incorrect patch version: %d", conf.MinVersion.Patch)
	}
}

func testServers(t *testing.T) {
	for server, pGroup := range conf.Servers {
		t.Logf("server: %s, pGroup: %v", server, pGroup)
	}
}
