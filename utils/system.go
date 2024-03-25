package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func updateSystem() error {
	fmt.Println("System is updating...")
	updateCommand := exec.Command("sudo", "apt", "update", "-y")
	updateCommand.Stdout = os.Stdout
	updateCommand.Stderr = os.Stderr
	err := updateCommand.Run()
	if err != nil {
		return err
	}
	return nil
}

func createDirectory(path string) error {
	fmt.Printf("Creates directory...")
	mkdirCommand := exec.Command("sudo", "mkdir", "-p", path)
	mkdirCommand.Stdout = os.Stdout
	mkdirCommand.Stderr = os.Stderr
	return mkdirCommand.Run()
}
