package tmux

import (
	"fmt"
	"os"
	"os/exec"
)

type Terminal struct {
	Command string `yaml:"command"`
}

type Window struct {
	//values: horizontal or vertical
	Split     string     `yaml:"split,omitempty"`
	Name      string     `yaml:"name"`
	Terminals []Terminal `yaml:"terminals"`
}

type Ash struct {
	Path          string   `yaml:"path"`
	SessionName   string   `yaml:"sessionName"`
	DefaultWindow string   `yaml:"defaultWindow"`
	Windows       []Window `yaml:"windows"`
}

func NewSession(ash Ash) {
	cmd := exec.Command(
		"tmux",
		"new-session",
		"-d",
		fmt.Sprintf("-s %s", ash.SessionName),
		"-c",
		ash.Path,
	)
	fmt.Println(cmd.Args)
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
	cmd := exec.Command(
		"tmux",
		"rename-window",
		fmt.Sprintf("-t %s", target),
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
		fmt.Sprintf("-t %s", ash.SessionName),
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
	cmd := exec.Command(
		"tmux",
		"send-keys",
		fmt.Sprintf("-t %s", target),
		command,
		"C-m",
	)
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
		fmt.Sprintf("-t %s", target),
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
		fmt.Sprintf("-t %s", ash.SessionName),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to attach session: %s\n", err)
	}
}

func HasSession(sessionName string) bool {
	cmd := exec.Command(
		"tmux",
		"has-session",
		fmt.Sprintf("-t= %s", sessionName),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		return false
	}
	return true
}

func switchSession(sessionName string) {
	cmd := exec.Command(
		"tmux",
		"switch-client",
		fmt.Sprintf("-t %s", sessionName),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to switch to session: %s\n", err)
	}
}

func ChangeSession(ash Ash) {
	tmuxEnv, tmuxEnvExist := os.LookupEnv("TMUX")
	fmt.Printf("$TMUX=%s exist:%v\n", tmuxEnv, tmuxEnvExist)
	if tmuxEnvExist {
		switchSession(ash.SessionName)
		return
	}
	Attach(ash)
}
