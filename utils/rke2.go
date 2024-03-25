package utils

import (
	"html/template"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/vmindtech/vke-agent/models"
)

func RKE2Install(version string, rke2AgentType string) error {
	logrus.Info("RKE2 Install...")
	curlCommand := "curl -sfL https://get.rke2.io | INSTALL_RKE2_VERSION=" + version + " INSTALL_RKE2_TYPE=" + rke2AgentType + " sh -"
	rke2InstallCommand := exec.Command("sh", "-c", curlCommand)
	rke2InstallCommand.Stdout = os.Stdout
	rke2InstallCommand.Stderr = os.Stderr
	return rke2InstallCommand.Run()
}

func RKE2ServiceEnable(rke2AgentType string) error {
	logrus.Info("RKE2 Enabled...")
	serviceName := "rke2-server"
	if rke2AgentType == "agent" {
		serviceName = "rke2-agent"
	}

	rke2ServiceEnableCommand := exec.Command("sudo", "systemctl", "enable", serviceName)
	rke2ServiceEnableCommand.Stdout = os.Stdout
	rke2ServiceEnableCommand.Stderr = os.Stderr
	return rke2ServiceEnableCommand.Run()
}

func RKE2ServiceStart(rke2AgentType string) error {
	logrus.Info("RKE2 started...")
	serviceName := "rke2-server"
	if rke2AgentType == "agent" {
		serviceName = "rke2-agent"
	}

	rke2ServiceStartCommand := exec.Command("sudo", "systemctl", "start", serviceName)
	rke2ServiceStartCommand.Stdout = os.Stdout
	rke2ServiceStartCommand.Stderr = os.Stderr
	return rke2ServiceStartCommand.Run()
}

func RKE2Config(initialize bool, serverAddress string, rke2AgentType string, rke2Token string, TlsSan string) error {
	logrus.Info("RKE2 config creating...")

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Error("Error getting hostname:", err)
		return err
	}

	cluster := []models.InitMaster{
		{
			NodeName:      hostname,
			Token:         rke2Token,
			TlsSan:        TlsSan,
			Initialize:    initialize,
			ServerAddress: serverAddress,
			Rke2AgentType: rke2AgentType,
		},
	}

	var yamlFile = "templates/config.yaml"
	yaml, err := template.New(yamlFile).ParseFiles(yamlFile)
	if err != nil {
		logrus.Error("Error parsing YAML file:", err)
		return err
	}

	f, err := os.Create("/etc/rancher/rke2/config.yaml")
	if err != nil {
		logrus.Error("Error creating config.yaml file:", err)
		return err
	}
	defer f.Close()

	err = yaml.Execute(f, cluster)
	if err != nil {
		logrus.Error("Error executing YAML template:", err)
		return err
	}

	logrus.Info("RKE2 config created.")
	return nil
}
