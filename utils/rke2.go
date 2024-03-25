package utils

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"

	"github.com/vmindtech/vke-agent/models"
)

func rke2Install(version string, rke2AgentType string) error {
	fmt.Println("RKE2 Install...")
	curlCommand := "curl -sfL https://get.rke2.io | INSTALL_RKE2_VERSION=" + version + " INSTALL_RKE2_TYPE=" + rke2AgentType + " sh -"
	rke2InstallCommand := exec.Command("sh", "-c", curlCommand)
	rke2InstallCommand.Stdout = os.Stdout
	rke2InstallCommand.Stderr = os.Stderr
	return rke2InstallCommand.Run()
}

func rke2ServiceStart(rke2AgentType string) error {
	fmt.Println("RKE2 started...")
	if rke2AgentType == "agent" {
		rke2ServiceStartCommand := exec.Command("sudo", "systemctl", "start", "rke2-agent")
		rke2ServiceStartCommand.Stdout = os.Stdout
		rke2ServiceStartCommand.Stderr = os.Stderr
		return rke2ServiceStartCommand.Run()
	} else {
		rke2ServiceStartCommand := exec.Command("sudo", "systemctl", "start", "rke2-server")
		rke2ServiceStartCommand.Stdout = os.Stdout
		rke2ServiceStartCommand.Stderr = os.Stderr
		return rke2ServiceStartCommand.Run()
	}
}

func rke2Config(initialize bool, serverAddress string, rke2AgentType string, rke2Token string, TlsSan string) error {
	fmt.Println("RKE2 config creating...")
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
	var yamlFile = "config.yaml"
	yaml, err := template.New(yamlFile).ParseFiles(yamlFile)
	f, err := os.Create("/etc/rancher/rke2/config.yaml")
	if err != nil {
		fmt.Println("Error creating config.yaml file:", err)
		return err
	}
	err = yaml.Execute(f, cluster)
	f.Close()
	return err
}
