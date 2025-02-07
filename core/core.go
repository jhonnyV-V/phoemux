package core

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/jhonnyV-V/phoemux/tmux"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goccy/go-yaml"
)

var (
	Quit bool
)

var (
	OpenEditor = true
)

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			//ignore case
			return false
		} else {
			fmt.Printf("failed to get file %s: %s\n", path, err)
			return false
		}
	}
	return true
}

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

func GetConfigPath() string {
	userConfigPath, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("failed to get config dir: %s\n", err)
		os.Exit(2)
	}

	return userConfigPath + "/phoemux"
}

func CreateConfigDir() string {
	phoemuxConfigPath := GetConfigPath()
	_, err := os.Stat(phoemuxConfigPath)
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

	return phoemuxConfigPath
}

func Create(phoemuxConfigPath, pwd, alias string) {

	if alias == "" {
		fmt.Printf("create command expects an alias\n")
		return
	}

	filePath := fmt.Sprintf(
		"%s/%s.yaml",
		phoemuxConfigPath,
		alias,
	)

	exist := fileExist(filePath)

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

	if !OpenEditor {
		return
	}

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

func Edit(phoemuxConfigPath, alias string) {
	if alias == "" {
		fmt.Printf("create command expects an alias\n")
		return
	}

	filePath := fmt.Sprintf(
		"%s/%s.yaml",
		phoemuxConfigPath,
		alias,
	)

	if !fileExist(filePath) {
		fmt.Printf("Ash %s does not exist\n", alias)
		return
	}

	editor := getEditor()
	cmd := exec.Command("sh", "-c", editor+" "+filePath)
	cmd.Env = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		fmt.Printf("failed to open editor: %s\n", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error while editing the file: %s\n", err)
	}
}

func getListOfItems(ashes []fs.DirEntry) []list.Item {
	items := []list.Item{}
	for _, ash := range ashes {
		if !strings.Contains(ash.Name(), ".yaml") {
			continue
		}
		name, _, _ := strings.Cut(ash.Name(), ".yaml")
		//TODO: display path inside file
		items = append(items, item(name))
	}
	return items
}

func ListAshes(phoemuxConfigPath string) {
	ashes, err := os.ReadDir(phoemuxConfigPath)
	if err != nil {
		fmt.Printf("Failed to read directory: %s\n", err)
	}

	var items []list.Item = getListOfItems(ashes)

	const defaultWidth = 20
	listKeys := newListKeyMap()

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Ashes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.openSelection,
			listKeys.editSelection,
			listKeys.deleteSelection,
		}
	}

	m := model{list: l}

	program := tea.NewProgram(m)

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if Quit {
		os.Exit(0)
	}

	switch Choice.Type {
	/* TODO: move delete and edit to the list.go file, and DO NOT QUIT on delete or edit
		consider using tea.ExecProcess to edit the file
	*/
	case "open":
		recreateFromAshes(phoemuxConfigPath, Choice.Target)
	case "delete":
		Delete(phoemuxConfigPath, Choice.Target)
	case "edit":
		Edit(phoemuxConfigPath, Choice.Target)
	}
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

func OpenFromCache(phoemuxConfigPath string) {
	cachePath := fmt.Sprintf(
		"%s/cache",
		phoemuxConfigPath,
	)

	if !fileExist(cachePath) {
		fmt.Printf("failed to get cache file\n")
		os.Exit(5)
	}

	file, err := os.ReadFile(cachePath)
	if err != nil {
		fmt.Printf("failed to read from cache %v\n", err)
		os.Exit(6)
	}

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

	if tmux.HasSession(ash.SessionName) {
		tmux.ChangeSession(ash)
		return
	}

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
	tmux.ChangeSession(ash)
}

func Delete(phoemuxConfigPath, alias string) {
	if alias == "" {
		fmt.Printf("delete command expects an alias\n")
		return
	}
	exist := ashExist(phoemuxConfigPath, alias)
	if !exist {
		fmt.Printf("Ash does not exist\n")
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

func Open(phoemuxConfigPath, alias string) {
	exist := ashExist(phoemuxConfigPath, alias)
	if exist {
		fmt.Printf("creating session\n")
		recreateFromAshes(phoemuxConfigPath, alias)
	} else {
		fmt.Printf("ash not found, can not create session\n")
	}
}

func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "nano"
	}

	return editor
}
