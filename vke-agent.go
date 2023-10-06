package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type InitMaster struct {
	NodeName      string
	Token         string
	TlsSan        string
	Initialize    bool
	ServerAddress string
	Rke2AgentType string
}

// sistem güncelleme
func updateSystem() error {
	fmt.Println(InfoColor, "System is updating...")
	updateCommand := exec.Command("sudo", "apt", "update", "-y")
	updateCommand.Stdout = os.Stdout
	updateCommand.Stderr = os.Stderr
	return updateCommand.Run()
}

// rke2 dizin oluşturma
func createDirectory(path string) error {
	fmt.Printf(InfoColor, "'%s' creates directory...\n", path)
	mkdirCommand := exec.Command("sudo", "mkdir", "-p", path)
	mkdirCommand.Stdout = os.Stdout
	mkdirCommand.Stderr = os.Stderr
	return mkdirCommand.Run()
}

// RKE2 yükleme
func rke2Install(version string, rke2AgentType string) error {
	fmt.Println(InfoColor, "RKE2 Install...")
	curlCommand := "curl -sfL https://get.rke2.io | INSTALL_RKE2_VERSION=" + version + " INSTALL_RKE2_TYPE=" + rke2AgentType + " sh -"
	rke2InstallCommand := exec.Command("sh", "-c", curlCommand)
	rke2InstallCommand.Stdout = os.Stdout
	rke2InstallCommand.Stderr = os.Stderr
	return rke2InstallCommand.Run()
}

func rke2ServiceStart(rke2AgentType string) error {
	fmt.Println(InfoColor, "RKE2 started...")
	if rke2AgentType == "server" {
		rke2ServiceStartCommand := exec.Command("sudo", "systemctl", "start", "rke2-server")
		rke2ServiceStartCommand.Stdout = os.Stdout
		rke2ServiceStartCommand.Stderr = os.Stderr
		return rke2ServiceStartCommand.Run()
	} else {
		rke2ServiceStartCommand := exec.Command("sudo", "systemctl", "start", "rke2-agent")
		rke2ServiceStartCommand.Stdout = os.Stdout
		rke2ServiceStartCommand.Stderr = os.Stderr
		return rke2ServiceStartCommand.Run()
	}
}
func rke2ServiceEnable(rke2AgentType string) error {
	fmt.Println(InfoColor, "RKE2 Enabled...")
	if rke2AgentType == "server" {
		rke2ServiceEnableCommand := exec.Command("sudo", "systemctl", "enable", "rke2-server")
		rke2ServiceEnableCommand.Stdout = os.Stdout
		rke2ServiceEnableCommand.Stderr = os.Stderr
		return rke2ServiceEnableCommand.Run()
	} else {
		rke2ServiceEnableCommand := exec.Command("sudo", "systemctl", "enable", "rke2-agent")
		rke2ServiceEnableCommand.Stdout = os.Stdout
		rke2ServiceEnableCommand.Stderr = os.Stderr
		return rke2ServiceEnableCommand.Run()
	}
}

// RKE2 Config oluşturma
func rke2Config(initialize bool, serverAddress string, rke2AgentType string, rke2Token string, TlsSan string) error {
	fmt.Println(InfoColor, "RKE2 config creating...")
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cluster := []InitMaster{
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
		// handle error
	}
	err = yaml.Execute(f, cluster)
	f.Close()
	return err
}

func main() {
	serverAddress := flag.String("serverAddress", "", "Server Address")
	kubeversion := flag.String("kubeversion", "", "Kube Version")
	//clusterID := flag.String("clusterID", "", "Cluster ID")
	tlsSan := flag.String("tlsSan", "", "TLS San")
	initialize := flag.Bool("initialize", false, "Initialize")
	rke2Token := flag.String("rke2Token", "", "RKE2 Token")
	rke2AgentType := flag.String("rke2AgentType", "", "Type")

	flag.Parse()

	if err := updateSystem(); err != nil {
		fmt.Println(ErrorColor, "System update error:", err)
		return
	}

	if err := createDirectory("/etc/rancher/rke2"); err != nil {
		fmt.Println(ErrorColor, "Indexing error:", err)
		return
	}
	if err := rke2Config(*initialize, *serverAddress, *rke2AgentType, *rke2Token, *tlsSan); err != nil {
		fmt.Println(ErrorColor, "Config creation error:", err)
		return
	}

	if err := rke2Install(*kubeversion, *rke2AgentType); err != nil {
		fmt.Println(ErrorColor, "RKE2 installation error:", err)
		return
	}

	if err := rke2ServiceEnable(*rke2AgentType); err != nil {
		fmt.Println(ErrorColor, "Service enabled error:", err)
		return
	}
	if err := rke2ServiceStart(*rke2AgentType); err != nil {
		fmt.Println(ErrorColor, "Service initialization error:", err)
		return
	}

	fmt.Println(InfoColor, "Process completed.")
}
