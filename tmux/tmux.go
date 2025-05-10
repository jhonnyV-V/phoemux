package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
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
		"-s", ash.SessionName,
		"-d",
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
		fmt.Sprintf("-t=%s", target),
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
		fmt.Sprintf("-t=%s", ash.SessionName),
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
		fmt.Sprintf("-t=%s", target),
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
		fmt.Sprintf("-t=%s", target),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to select window: %s\n", err)
	}
}

func switchSession(sessionName string) {
	cmd := exec.Command(
		"tmux",
		"switch-client",
		fmt.Sprintf("-t=%s", sessionName),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to switch to session: %s\n", err)
	}
}

func Attach(ash Ash) {
	var cmd *exec.Cmd
	if ash.SessionName != "" {
		cmd = exec.Command(
			"tmux",
			"attach-session",
			fmt.Sprintf("-t=%s", ash.SessionName),
		)
	} else {
		cmd = exec.Command(
			"tmux",
			"attach-session",
		)
	}

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
		fmt.Sprintf("-t=%s", sessionName),
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

func IsInsideTmux() bool {
	_, tmuxEnvExist := os.LookupEnv("TMUX")
	return tmuxEnvExist
}

func ChangeSession(ash Ash) {
	tmuxEnvExist := IsInsideTmux()
	if tmuxEnvExist {
		switchSession(ash.SessionName)
		return
	}
	Attach(ash)
}

func GetCurrentSessionName() string {
	cmd := exec.Command(
		"tmux",
		"display-message",
		"-p",
		"#S",
	)

	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(string(out), "\n", "")
}

func GetListOfSessions() []string {
	sessions := []string{}
	cmd := exec.Command(
		"tmux",
		"list-sessions",
		"-F",
		"#{session_name}",
	)

	out, err := cmd.Output()
	if err != nil {
		return sessions
	}
	sessions = strings.Split(string(out), "\n")
	return filter(sessions, func(s string) bool {
		return s != ""
	})
}

func filter[T any](slice []T, f func(T) bool) []T {
	for i, value := range slice {
		if !f(value) {
			result := slices.Clone(slice[:i])
			for i++; i < len(slice); i++ {
				value = slice[i]
				if f(value) {
					result = append(result, value)
				}
			}
			return slices.Clip(result)
		}
	}
	return slice
}

func GetOthersSessions() []string {
	sessions := GetListOfSessions()
	currentSession := GetCurrentSessionName()
	sessions = filter(sessions, func(s string) bool {
		return strings.TrimSpace(s) != strings.TrimSpace(currentSession)
	})

	return sessions
}

func GetOtherSession() string {
	sessions := GetListOfSessions()
	currentSession := GetCurrentSessionName()
	sessions = filter(sessions, func(s string) bool {
		return strings.TrimSpace(s) != strings.TrimSpace(currentSession)
	})
	if len(sessions) == 0 {
		return ""
	}
	return sessions[0]
}

func Kill(sessionName string) {
	var cmd *exec.Cmd

	if sessionName != "" {
		cmd = exec.Command(
			"tmux",
			"kill-session",
			"-t", sessionName,
		)
	} else {
		cmd = exec.Command(
			"tmux",
			"kill-session",
		)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to kill session: %s\n", err)
	}
}
