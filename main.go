package main

import (
	"flag"
	"fmt"
	"github.com/jhonnyV-V/phoemux/tmux"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goccy/go-yaml"
)

var (
	Quit bool
)

//TODO: add last or restore command to reopen the last ash

func getDefault(path, alias string) string {
	return fmt.Sprintf(`path: "%s"
sessionName: "%s"
defaultWindow: code
windows:
- name: code
  terminals:
  - command: echo "do something here"
- name: servers
  terminals:
  - command: ls`,
		path,
		alias,
	)
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get pwd: %s\n", err)
	}

	userConfigPath, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("failed to get config dir: %s\n", err)
		os.Exit(2)
	}

	phoemuxConfigPath := userConfigPath + "/phoemux"

	_, err = os.Stat(phoemuxConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(phoemuxConfigPath, 0766)
			if err != nil {
				fmt.Printf("failed to create phoenix dir: %s\n", err)
				os.Exit(3)
			}
		} else {
			fmt.Printf("failed to get phoenix dir: %s\n", err)
			os.Exit(4)
		}
	}

	flag.Parse()

	command := flag.Arg(0)
	//fmt.Printf("args %#v \n", flag.Args())

	switch command {
	case "create":
		create(
			phoemuxConfigPath,
			dir,
		)

	case "edit":
		edit(phoemuxConfigPath)

	case "list":
		listAshes(phoemuxConfigPath)

	case "delete":
		delete(phoemuxConfigPath)

	case "last":
		openFromCache(phoemuxConfigPath)

	case "":
		fmt.Printf("empty command\n")

	default:
		fmt.Printf("unkown command maybe rise from the ashes\n")
		exist := ashExist(phoemuxConfigPath, command)
		fmt.Printf("ash exist\n")
		if exist {
			fmt.Printf("creating session\n")
			recreateFromAshes(phoemuxConfigPath, command)
		}
	}
}

func create(phoemuxConfigPath, pwd string) {
	alias := flag.Arg(1)
	exist := true

	if alias == "" {
		fmt.Printf("create command expects an alias\n")
		return
	}

	filePath := fmt.Sprintf(
		"%s/%s.yaml",
		phoemuxConfigPath,
		alias,
	)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			//ignore case
			exist = false
		} else {
			fmt.Printf("failed to get ash for %s: %s\n", alias, err)
			return
		}
	}

	if exist {
		fmt.Printf("ash for %s already exist\n", alias)
		fmt.Printf("if you want to edit it use the edit command\n")
		return
	}

	config, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create ash: %s\n", err)
		return
	}

	example := getDefault(pwd, alias)

	_, err = config.Write([]byte(example))
	if err != nil {
		fmt.Printf("Failed write ash: %s\n", err)
		return
	}
	config.Close()

	editor := getEditor()
	cmd := exec.Command("sh", "-c", editor+" "+filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Start()
	if err != nil {
		fmt.Printf("failed to open editor: %s\n", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error while editing the file: %s\n", err)
	}
}

func edit(phoemuxConfigPath string) {
	alias := flag.Arg(1)

	if alias == "" {
		fmt.Printf("create command expects an alias\n")
		return
	}

	filePath := fmt.Sprintf(
		"%s/%s.yaml",
		phoemuxConfigPath,
		alias,
	)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			//ignore case
			fmt.Printf("ash for %s does not exist\n", alias)
			return
		} else {
			fmt.Printf("failed to get ash for %s: %s\n", alias, err)
			return
		}
	}

	editor := getEditor()
	cmd := exec.Command("sh", "-c", editor+" "+filePath)
	cmd.Env = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Start()
	if err != nil {
		fmt.Printf("failed to open editor: %s\n", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error while editing the file: %s\n", err)
	}
}

func listAshes(phoemuxConfigPath string) {
	ashes, err := os.ReadDir(phoemuxConfigPath)
	if err != nil {
		fmt.Printf("Failed to read directory: %s\n", err)
	}

	items := []list.Item{}

	for _, ash := range ashes {
		if !strings.Contains(ash.Name(), ".yaml") {
			continue
		}
		name, _, _ := strings.Cut(ash.Name(), ".yaml")
		//TODO: display path inside file
		items = append(items, item(name))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Ashes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if Quit {
		os.Exit(0)
	}

	recreateFromAshes(phoemuxConfigPath, Choice)
}

func writeToCache(phoemuxConfigPath, alias string) {
	cachePath := fmt.Sprintf(
		"%s/cache",
		phoemuxConfigPath,
	)

	err := os.WriteFile(cachePath, []byte(alias), 0766)
	if err != nil {
		fmt.Printf("Failed to write to cache: %s\n", err)
		return
	}
}

func openFromCache(phoemuxConfigPath string) {
	cachePath := fmt.Sprintf(
		"%s/cache",
		phoemuxConfigPath,
	)

	_, err := os.Stat(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("failed to read from cache: %s\n", err)
			os.Exit(5)
		} else {
			fmt.Printf("failed to get cache file: %s\n", err)
			os.Exit(6)
		}
	}

	file, err := os.ReadFile(cachePath)

	recreateFromAshes(phoemuxConfigPath, string(file))
}

func recreateFromAshes(phoemuxConfigPath, alias string) {
	var ash tmux.Ash

	filePath := fmt.Sprintf(
		"%s/%s.yaml",
		phoemuxConfigPath,
		alias,
	)

	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read ash: %s\n", err)
		return
	}

	err = yaml.Unmarshal(file, &ash)
	if err != nil {
		fmt.Printf("Failed to unmarshall ash: %s\n", err)
		return
	}

	writeToCache(phoemuxConfigPath, alias)

	tmux.NewSession(ash)
	for i, window := range ash.Windows {
		if i == 0 {
			tmux.RenameWindow(ash, "0", window.Name)
		} else {
			tmux.NewWindow(ash, window)
		}

		tmux.RunCommand(
			ash.SessionName,
			window.Name,
			window.Terminals[0].Command,
		)
	}

	tmux.SetWindows(ash)
	tmux.Attach(ash)
}

func delete(phoemuxConfigPath string) {
	alias := flag.Arg(1)

	if alias == "" {
		fmt.Printf("delete command expects an alias\n")
		return
	}
	exist := ashExist(phoemuxConfigPath, alias)
	if !exist {
		fmt.Printf("ash does not exist\n")
		return
	}

	os.Remove(
		phoemuxConfigPath + "/" + alias + ".yaml",
	)
}

func ashExist(phoemuxConfigPath, alias string) bool {
	ashes, err := os.ReadDir(phoemuxConfigPath)
	if err != nil {
		fmt.Printf("Failed to read directory: %s\n", err)
	}

	for _, ash := range ashes {
		if !strings.Contains(ash.Name(), ".yaml") {
			continue
		}
		name, _, _ := strings.Cut(ash.Name(), ".yaml")
		if alias == name {
			return true
		}
	}
	return false
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "nano"
	}

	return editor
}
