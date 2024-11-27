package core

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	//before tests
	configPath := GetConfigPath()
	oldConfigPath := path.Join(configPath, "..", "old-phoemux")
	fmt.Printf("move config folder to %s\n", oldConfigPath)
	err := os.Rename(configPath, oldConfigPath)
	if err != nil {
		fmt.Printf("failed to move config directory to %s with error %s\n", oldConfigPath, err)
		os.Exit(1)
	}
	CreateConfigDir()
	OpenEditor = false

	exitVal := m.Run()
	//cleanup
	err = os.RemoveAll(configPath)
	if err != nil {
		fmt.Printf("failed remove test config directory with error %s\n", err)
		os.Exit(1)
	}
	err = os.Rename(oldConfigPath, configPath)
	if err != nil {
		fmt.Printf("failed to move old config directory to %s with error %s\n", configPath, err)
		os.Exit(1)
	}
	os.Exit(exitVal)
}

func TestCreate(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get pwd: %s\n", err)
		os.Exit(1)
	}
	phoemuxConfigPath := CreateConfigDir()
	Create(phoemuxConfigPath, pwd, "phoemux")
	created := ashExist(phoemuxConfigPath, "phoemux")
	if !created {
		t.Fatal("failed to create file")
	}
}

func getList(ashes []fs.DirEntry) []string {
	items := []string{}
	for _, ash := range ashes {
		if !strings.Contains(ash.Name(), ".yaml") {
			continue
		}
		name, _, _ := strings.Cut(ash.Name(), ".yaml")
		//TODO: display path inside file
		items = append(items, name)
	}
	return items
}

func TestList(t *testing.T) {
	path := GetConfigPath()
	ashes, err := os.ReadDir(path)
	if err != nil {
		t.Fatalf("failed to read dir %s with error %s\n", path, err)
	}
	items := getList(ashes)
	if len(items) == 0 {
		t.Fatal("no elements in config folder")
	}
	if items[0] != "phoemux" {
		t.Fatalf("unexpected value %s\n", items[0])
	}
}

func TestDelete(t *testing.T) {
	phoemuxConfigPath := GetConfigPath()

	Delete(phoemuxConfigPath, "phoemux")

	deleted := ashExist(phoemuxConfigPath, "phoemux")
	if deleted {
		t.Fatal("failed to delete file")
	}
}
