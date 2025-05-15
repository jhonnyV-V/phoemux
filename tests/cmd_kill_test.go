package tests

import (
	"testing"

	"github.com/jhonnyV-V/phoemux/core"
	"github.com/jhonnyV-V/phoemux/tmux"
)

func TestKill(t *testing.T) {
	phoemuxConfigPath := core.GetConfigPath()
	target := "vc"

	core.Open(phoemuxConfigPath, target)

	if !tmux.HasSession(target) {
		t.Fatalf("Failed to created the session\n")
	}

	tmux.Kill(target)

	if tmux.HasSession(target) {
		t.Fatalf("Failed to kill target: %s\n", target)
	}
}
