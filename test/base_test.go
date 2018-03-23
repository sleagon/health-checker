package test

import (
	"github.com/sleagon/health-checker/base"
	"testing"
)

func TestLoadConfigOk(t *testing.T) {
	base.LoadConfig("config.ok.json")
	cfg := base.GetConfig()
	if cfg.ProjectName != "HealthCheckerOk" {
		t.Error("Failed to load legal config file.")
	}
}

func TestLoadConfigErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Panic should be fired when using a broken config file.")
			return
		}
	}()
	base.LoadConfig("test/config.err.json")
	cfg := base.GetConfig()
	if cfg.ProjectName != "HealthCheckerOk" {
		t.Error("Failed to load legal config file.")
	}
}
