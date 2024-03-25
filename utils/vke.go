package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/vmindtech/vke-agent/models"
	"github.com/vmindtech/vke-agent/request"
	"gopkg.in/yaml.v2"
)

func PushRKE2Config(initialize bool, rke2AgentType, serverAddress, clusterName, ClusterUUID, VKEAPIEndpoint, VKEAPIAuthToken string) error {
	_, err := os.Stat("/etc/rancher/rke2/rke2.yaml")
	if os.IsNotExist(err) {
		fmt.Println("RKE2 config file not found")
		return fmt.Errorf("RKE2 config file not found")
	}

	if !initialize && rke2AgentType != "server" && serverAddress == "" && clusterName == "" && ClusterUUID == "" && VKEAPIEndpoint == "" && VKEAPIAuthToken == "" {
		fmt.Printf("RKE2 config insufficient parameters")
		return fmt.Errorf("RKE2 config insufficient parameters")
	}

	fmt.Println("RKE2 config pushing...")
	data, err := os.ReadFile("/etc/rancher/rke2/rke2.yaml")
	if err != nil {
		fmt.Println("Config reading error:", err)
		return err
	}

	var kubeconfig models.KubeConfig
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

	sendKubeConfigRequest := request.SendKubeConfigRequest{
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
