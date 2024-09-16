package main

import (
	"encoding/json"
	"flag"
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

func getDefault(path string) Ash {
	return Ash{
		Path:          path,
		SessionName:   "my_project",
		DefaultWindow: "code",
		Windows: []Window{
			{
				Name: "code",
				Terminals: []Terminal{
					{
						Command: "echo \"do something here \"",
					},
				},
			},
		},
	}
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get pwd: %s\n", err)
	}
	fmt.Printf("PWD: %s\n", dir)

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
	fmt.Printf("args %#v \n", flag.Args())

	switch command {
	case "create":
		create(
			phoemuxConfigPath,
			dir,
		)

	case "edit":
		fmt.Printf("edit command\n")

	case "list":
		fmt.Printf("list command\n")

	case "delete":
		fmt.Printf("delete command\n")

	case "":
		fmt.Printf("empty command\n")

	default:
		fmt.Printf("unkown command maybe rise from the ashes\n")
	}
}

func create(phoemuxConfigPath, pwd string) {
	fmt.Printf("create command\n")
	alias := flag.Arg(1)
	exist := true

	if alias == "" {
		fmt.Printf("create command expects an alias\n")
		return
	}

	filePath := fmt.Sprintf(
		"%s/%s.json",
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

	example := getDefault(pwd)

	data, err := json.Marshal(example)
	if err != nil {
		fmt.Printf("Failed to marshall ash: %s\n", err)
		return
	}
	_, err = config.Write(data)
	if err != nil {
		fmt.Printf("Failed write ash: %s\n", err)
		return
	}
	config.Close()

	//TODO: open in $EDITOR
	cmd := exec.Command("sh", "-c", "$EDITOR "+filePath)
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
