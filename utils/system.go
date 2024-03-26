package utils

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func init() {
	// Log seviyesini ayarla, eğer gerekliyse
	logrus.SetLevel(logrus.InfoLevel)
}

func UpdateSystem() error {
	logrus.Info("System is updating...")
	updateCommand := exec.Command("sudo", "apt", "update", "-y")
	updateCommand.Stdout = os.Stdout
	updateCommand.Stderr = os.Stderr
	err := updateCommand.Run()
	if err != nil {
		logrus.Error("System update error:", err)
		return err
	}
	logrus.Info("System update completed.")
	return nil
}

func CreateDirectory(path string) error {
	logrus.Infof("Creating directory: %s", path)
	mkdirCommand := exec.Command("sudo", "mkdir", "-p", path)
	mkdirCommand.Stdout = os.Stdout
	mkdirCommand.Stderr = os.Stderr
	err := mkdirCommand.Run()
	if err != nil {
		logrus.Error("Directory creation error:", err)
		return err
	}
	logrus.Infof("Directory created: %s", path)
	return nil
}
