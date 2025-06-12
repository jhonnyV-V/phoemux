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

func GetListOfWindows(sessionName string) ([]string, string) {
	windows := []string{}
	active := ""
	cmd := exec.Command(
		"tmux",
		"list-windows",
		"-t", sessionName,
		"-F",
		"#{window_name}active=#{window_active}",
	)

	out, err := cmd.Output()
	if err != nil {
		return windows, active
	}
	windows = strings.Split(string(out), "\n")

	windows = filter(windows, func(s string) bool {
		return s != ""
	})

	for i, v := range windows {
		result := strings.Split(v, "active=")
		if result[1] == "1" {
			active = result[0]
		}
		windows[i] = result[0]
	}
	return windows, active
}

func GetListOfPanes(sessionName string) []string {
	cmd := exec.Command(
		"tmux",
		"list-panes",
		"-a",
		"-F",
		"#{pane_id} #{pane_current_command} #{session_name}",
	)

	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("Failed to get list of panes: %s\n", err)
		return []string{}
	}

	panes := strings.Split(string(out), "\n")

	panes = filter(panes, func(s string) bool {
		if s == "" {
			return false
		}

		return strings.Contains(s, sessionName)
	})

	return panes
}

func SendCommandToPane(paneId string, commands []string) {
	args := []string{"send-keys", "-t", paneId}

	for _, command := range commands {
		args = append(args, command)
	}
	cmd := exec.Command(
		"tmux",
		args...,
	)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to send command to pane: %s\n", err)
	}
}

func killAllProceessInSession(sessionName string) {
	panes := GetListOfPanes(sessionName)
	slices.Reverse(panes)

	for _, pane := range panes {
		paneData := strings.Split(pane, " ")
		paneId := paneData[0]
		paneProc := strings.ToLower(paneData[1])
		cmd := []string{}
		if paneProc == "vim" || paneProc == "vi" || paneProc == "nvim" {
			cmd = append(cmd, "Escape")
			cmd = append(cmd, ":qa")
			cmd = append(cmd, "Enter")
		} else if paneProc == "emacs" {
			cmd = append(cmd, "C-x")
			cmd = append(cmd, "C-c")
		} else if paneProc == "man" || paneProc == "less" {
			cmd = append(cmd, "q")
		} else if paneProc == "bash" || paneProc == "zsh" || paneProc == "fish" {
			cmd = append(cmd, "C-c")
			cmd = append(cmd, "C-u")
			cmd = append(cmd, "space")
			cmd = append(cmd, "\"exit\"")
			cmd = append(cmd, "Enter")
		} else if paneProc == "ssh" || paneProc == "vagrant" {
			cmd = append(cmd, "Enter")
			cmd = append(cmd, "\"~.\"")
		} else if paneProc == "psql" || paneProc == "mysql" {
			cmd = append(cmd, "C-d")
		} else if paneProc == "go" && paneData[2] == "phoemux" {
			cmd = append(cmd, "")
		} else if paneProc == "phoemux" {
			cmd = append(cmd, "")
		}

		if len(cmd) == 0 {
			cmd = append(cmd, "C-c")
			cmd = append(cmd, "C-c")
			cmd = append(cmd, "C-c")
			cmd = append(cmd, "C-c")
			cmd = append(cmd, "C-c")
			cmd = append(cmd, "C-c")
			cmd = append(cmd, "C-c")
			cmd = append(cmd, "C-c")
		}

		SendCommandToPane(paneId, cmd)
	}
}

func Kill(sessionName string) {

	cmd := exec.Command(
		"tmux",
		"kill-session",
		"-t", sessionName,
	)

	killAllProceessInSession(sessionName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to kill session: %s\n", err)
	}
}
