package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

func updateSystem() error {
	fmt.Println("System is updating...")
	updateCommand := exec.Command("sudo", "apt", "update", "-y")
	updateCommand.Stdout = os.Stdout
	updateCommand.Stderr = os.Stderr
	return updateCommand.Run()
}

func createDirectory(path string) error {
	fmt.Printf("Creates directory...")
	mkdirCommand := exec.Command("sudo", "mkdir", "-p", path)
	mkdirCommand.Stdout = os.Stdout
	mkdirCommand.Stderr = os.Stderr
	return mkdirCommand.Run()
}

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
func rke2ServiceEnable(rke2AgentType string) error {
	fmt.Println("RKE2 Enabled...")
	if rke2AgentType == "agent" {
		rke2ServiceEnableCommand := exec.Command("sudo", "systemctl", "enable", "rke2-agent")
		rke2ServiceEnableCommand.Stdout = os.Stdout
		rke2ServiceEnableCommand.Stderr = os.Stderr
		return rke2ServiceEnableCommand.Run()
	} else {
		rke2ServiceEnableCommand := exec.Command("sudo", "systemctl", "enable", "rke2-server")
		rke2ServiceEnableCommand.Stdout = os.Stdout
		rke2ServiceEnableCommand.Stderr = os.Stderr
		return rke2ServiceEnableCommand.Run()
	}
}

func rke2Config(initialize bool, serverAddress string, rke2AgentType string, rke2Token string, TlsSan string) error {
	fmt.Println("RKE2 config creating...")
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
		fmt.Println("Error creating config.yaml file:", err)
		return err
	}
	err = yaml.Execute(f, cluster)
	f.Close()
	return err
}

func pushRKE2Config(initialize bool, rke2AgentType, serverAddress, clusterName, ClusterUUID, VKEAPIEndpoint, VKEAPIAuthToken string) error {
	_, err := os.Stat("./rke2-demo.yaml")
	if os.IsNotExist(err) {
		fmt.Println("RKE2 config file not found")
		return fmt.Errorf("RKE2 config file not found")
	}

	if !initialize && rke2AgentType != "server" && serverAddress == "" && clusterName == "" && ClusterUUID == "" && VKEAPIEndpoint == "" && VKEAPIAuthToken == "" {
		fmt.Printf("RKE2 config insufficient parameters")
		return fmt.Errorf("RKE2 config insufficient parameters")
	}

	fmt.Println("RKE2 config pushing...")
	data, err := os.ReadFile("./rke2-demo.yaml")
	if err != nil {
		fmt.Println("Config reading error:", err)
		return err
	}

	var kubeconfig KubeConfig
	err = yaml.Unmarshal([]byte(data), &kubeconfig)
	if err != nil {
		fmt.Println("Config unmarshal error:", err)
		return err
	}

	kubeconfig.Clusters[0].Cluster.Server = fmt.Sprintf("https://%s:6443", serverAddress)
	kubeconfig.Clusters[0].Name = clusterName

	kubeconfig.Contexts[0].Context.Cluster = clusterName
	kubeconfig.Contexts[0].Context.User = clusterName
	kubeconfig.Contexts[0].Name = clusterName
	kubeconfig.CurrentContext = clusterName

	kubeconfig.Users[0].Name = clusterName

	newKubeConfigYaml, err := yaml.Marshal(&kubeconfig)
	if err != nil {
		fmt.Println("Config marshal error:", err)
		return err
	}

	kubeConfigBase64 := base64.StdEncoding.EncodeToString(newKubeConfigYaml)

	sendKubeConfigRequest := SendKubeConfigRequest{
		ClusterID:  ClusterUUID,
		KubeConfig: kubeConfigBase64,
	}

	kubeConfigData, err := json.Marshal(sendKubeConfigRequest)
	if err != nil {
		fmt.Println("KubeConfig json marshal error:", err)
		return err
	}

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/kubeconfig", VKEAPIEndpoint), bytes.NewBuffer(kubeConfigData))
	if err != nil {
		fmt.Println("KubeConfig request error:", err)
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Auth-Token", VKEAPIAuthToken)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println("KubeConfig response error:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("KubeConfig response status code error:", resp.StatusCode)
		return err
	}

	return nil
}

func main() {
	serverAddress := flag.String("serverAddress", "", "Server Address")
	kubeversion := flag.String("kubeversion", "", "Kube Version")
	tlsSan := flag.String("tlsSan", "", "TLS San")
	initialize := flag.Bool("initialize", false, "Initialize")
	rke2Token := flag.String("rke2Token", "", "RKE2 Token")
	rke2AgentType := flag.String("rke2AgentType", "", "Type")
	rke2ClusterName := flag.String("rke2ClusterName", "", "Cluster Name")
	rke2ClusterUUID := flag.String("rke2ClusterUUID", "", "Cluster UUID")
	rke2AgentVKEAPIEndpoint := flag.String("rke2AgentVKEAPIEndpoint", "", "VKE API Endpoint")
	rke2AgentVKEAPIAuthToken := flag.String("rke2AgentVKEAPIAuthToken", "", "VKE API Auth Token")

	flag.Parse()

	if err := updateSystem(); err != nil {
		fmt.Println("System update error:", err)
		return
	}

	if err := createDirectory("/etc/rancher/rke2"); err != nil {
		fmt.Println("Indexing error:", err)
		return
	}
	if err := rke2Config(*initialize, *serverAddress, *rke2AgentType, *rke2Token, *tlsSan); err != nil {
		fmt.Println("Config creation error:", err)
		return
	}

	if err := rke2Install(*kubeversion, *rke2AgentType); err != nil {
		fmt.Println("RKE2 installation error:", err)
		return
	}

	if err := rke2ServiceEnable(*rke2AgentType); err != nil {
		fmt.Println("Service enabled error:", err)
		return
	}
	if err := rke2ServiceStart(*rke2AgentType); err != nil {
		fmt.Println("Service initialization error:", err)
		return
	}
	if err := pushRKE2Config(*initialize, *rke2AgentType, *serverAddress, *rke2ClusterName, *rke2ClusterUUID, *rke2AgentVKEAPIEndpoint, *rke2AgentVKEAPIAuthToken); err != nil {
		fmt.Println("Pushing RKE2 config error:", err)
		return
	}

	fmt.Println("Process completed.")
}
