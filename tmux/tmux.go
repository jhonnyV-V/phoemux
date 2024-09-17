package tmux

import (
	"fmt"
	"os"
	"os/exec"
)

type Terminal struct {
	Command string `json:"command"`
}

type Window struct {
	//values: horizontal or vertical
	Split     string     `json:"split,omitempty"`
	Name      string     `json:"name"`
	Terminals []Terminal `json:"terminals"`
}

type Ash struct {
	Path          string   `json:"path"`
	SessionName   string   `json:"sessionName"`
	DefaultWindow string   `json:"defaultWindow"`
	Windows       []Window `json:"windows"`
}

func NewSession(ash Ash) {
	cmd := exec.Command(
		"tmux",
		"new-session",
		"-d",
		"-s "+ash.SessionName,
		"-c",
		ash.Path,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to open session: %s\n", err)
	}
}

func RenameWindow(ash Ash, oldName, newName string) {
	target := fmt.Sprintf("%s:%s", ash.SessionName, oldName)
	fmt.Printf("rename target %s\n", target)
	cmd := exec.Command(
		"tmux",
		"rename-window",
		"-t "+target,
		newName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to rename window: %s\n", err)
	}
}

func NewWindow(ash Ash, window Window) {
	cmd := exec.Command(
		"tmux",
		"new-window",
		"-c",
		ash.Path,
		"-n",
		window.Name,
		"-t "+ash.SessionName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to create window %s\n", err)
	}
}

func RunCommand(sessionName, currentWindow, command string) {
	target := fmt.Sprintf("%s:%s.0", sessionName, currentWindow)
	cmd := exec.Command("tmux", "send-keys", "-t "+target, command, "C-m")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to run command %s: %s\n", command, err)
	}
}

func SetWindows(ash Ash) {
	target := fmt.Sprintf("%s:%s", ash.SessionName, ash.DefaultWindow)
	cmd := exec.Command(
		"tmux",
		"select-window",
		"-t "+target,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to select window: %s\n", err)
	}
}

func Attach(ash Ash) {
	cmd := exec.Command(
		"tmux",
		"attach-session",
		"-t "+ash.SessionName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to attach session: %s\n", err)
	}
}
