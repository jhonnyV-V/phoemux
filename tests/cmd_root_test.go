package tests

import (
	"testing"

	"github.com/jhonnyV-V/phoemux/core"
	"github.com/jhonnyV-V/phoemux/tmux"
)

func TestRoot(t *testing.T) {
	phoemuxConfigPath := core.GetConfigPath()
	target := "phoemux"
	shouldKill := tmux.GetCurrentSessionName() != target
	expectedActive := ""
	if shouldKill {
		expectedActive = "code"
	} else {
		expectedActive = "compiler"
	}

	expected := [3]string{
		"code",
		"compiler",
		"elevated",
	}

	core.Open(phoemuxConfigPath, target)

	if !tmux.HasSession(target) {
		t.Fatalf("Failed to created the session\n")
	}

	actualWindows, actualActive := tmux.GetListOfWindows(target)

	for i, v := range expected {
		if v != actualWindows[i] {
			t.Fatalf("Failed to created windows \nexpected %#v\nactual %#v at index %d\n", expected, actualWindows, i)
		}
	}

	if actualActive != expectedActive {
		t.Fatalf("Wrong active window \nexpected %s actual %s\n", expectedActive, actualActive)
	}

	if shouldKill {
		tmux.Kill(target)
	}
}
