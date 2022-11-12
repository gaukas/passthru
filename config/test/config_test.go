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
	if conf.Version.Major != 0 {
		t.Errorf("incorrect major version: %d", conf.Version.Major)
	}
	if conf.Version.Minor != 2 {
		t.Errorf("incorrect minor version: %d", conf.Version.Minor)
	}
	if conf.Version.Patch != 1 {
		t.Errorf("incorrect patch version: %d", conf.Version.Patch)
	}
}

func testServers(t *testing.T) {
	for server, pGroup := range conf.Servers {
		t.Logf("server: %s, pGroup: %v", server, pGroup)
	}
}
