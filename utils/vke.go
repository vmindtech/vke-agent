package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vmindtech/vke-agent/models"
	"github.com/vmindtech/vke-agent/request"
	"gopkg.in/yaml.v2"
)

func EditCanalHelmChartConfig() error {
	var yamlFile = "/var/lib/rancher/rke2/server/manifests/rke2-canal.yml"

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		logrus.Error("Error reading YAML file:", err)
		return err
	}

	var helmChartConfig models.HelmChartConfig
	err = yaml.Unmarshal(data, &helmChartConfig)
	if err != nil {
		logrus.Error("Error unmarshaling YAML file:", err)
		return err
	}

	if helmChartConfig.Spec.Set == nil {
		helmChartConfig.Spec.Set = make(map[string]interface{})
	}
	helmChartConfig.Spec.Set["calico.vethMTU"] = 1350

	newYamlData, err := yaml.Marshal(&helmChartConfig)
	if err != nil {
		logrus.Error("Error marshaling YAML:", err)
		return err
	}

	return os.WriteFile(yamlFile, newYamlData, 0644)
}

func PushRKE2Config(
	initialize bool,
	rke2AgentType,
	serverAddress,
	clusterName,
	ClusterUUID,
	VKEAPIEndpoint,
	VKEAPIAuthToken string,
) error {
	_, err := os.Stat("/etc/rancher/rke2/rke2.yaml")
	if os.IsNotExist(err) {
		logrus.Error("RKE2 config file not found")
		return err
	}

	if !initialize && rke2AgentType != "server" && serverAddress == "" && clusterName == "" && ClusterUUID == "" && VKEAPIEndpoint == "" && VKEAPIAuthToken == "" {
		logrus.Error("RKE2 config insufficient parameters")
		return err
	}

	logrus.Info("RKE2 config pushing...")
	data, err := os.ReadFile("/etc/rancher/rke2/rke2.yaml")
	if err != nil {
		logrus.Error("Config reading error:", err)
		return err
	}

	var kubeconfig models.KubeConfig
	err = yaml.Unmarshal([]byte(data), &kubeconfig)
	if err != nil {
		logrus.Error("Config unmarshal error:", err)
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
		logrus.Error("Config marshal error:", err)
		return err
	}

	kubeConfigBase64 := base64.StdEncoding.EncodeToString(newKubeConfigYaml)

	sendKubeConfigRequest := request.SendKubeConfigRequest{
		ClusterID:  ClusterUUID,
		KubeConfig: kubeConfigBase64,
	}

	kubeConfigData, err := json.Marshal(sendKubeConfigRequest)
	if err != nil {
		logrus.Error("KubeConfig json marshal error:", err)
		return err
	}

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/kubeconfig", VKEAPIEndpoint), bytes.NewBuffer(kubeConfigData))
	if err != nil {
		logrus.Error("KubeConfig request error:", err)
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Auth-Token", VKEAPIAuthToken)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		logrus.Error("KubeConfig response error:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logrus.Error("KubeConfig response status code error:", resp.StatusCode)
		return err
	}

	return nil
}
func DeployHelmCharts(
	ClusterUUID,
	RKE2ClusterProjectUUID,
	VkeCloudAuthURL,
	ApplicationCredentialID,
	ApplicationCredentialKey,
	CloudControllerManagerVersion,
	AutoScalerVersion,
	ClusterAgentVersion,
	RKE2AgentVKEAPIEndpoint,
	LoadBalancerFloatingNetworkID string,
) error {
	err := os.MkdirAll("/var/lib/rancher/rke2/server/manifests", 0755)
	if err != nil {
		logrus.Error("Error creating directory:", err)
		return err
	}

	var yamlFile = "k8s-helmchart-for-cloud-provider.yml"
	yaml, err := template.New(yamlFile).ParseFiles(yamlFile)
	if err != nil {
		logrus.Error("Error parsing YAML file:", err)
		return err
	}

	f, err := os.Create("/var/lib/rancher/rke2/server/manifests/k8s-helmchart-for-cloud-provider.yml")
	if err != nil {
		logrus.Error("Error creating k8s-helmchart-for-cloud-provider.yml file:", err)
		return err
	}
	defer f.Close()

	cluster := []models.InitMaster{
		{
			RKE2ClusterProjectUUID:        RKE2ClusterProjectUUID,
			RKE2ClusterUUID:               ClusterUUID,
			VkeCloudAuthURL:               VkeCloudAuthURL,
			ApplicationCredentialID:       ApplicationCredentialID,
			ApplicationCredentialKey:      ApplicationCredentialKey,
			ClusterAutoscalerVersion:      AutoScalerVersion,
			ClusterAgentVersion:           ClusterAgentVersion,
			CloudProviderVkeVersion:       CloudControllerManagerVersion,
			RKE2AgentVKEAPIEndpoint:       RKE2AgentVKEAPIEndpoint,
			LoadBalancerFloatingNetworkID: LoadBalancerFloatingNetworkID,
		},
	}
	err = yaml.Execute(f, cluster)
	if err != nil {
		logrus.Error("Error executing YAML template:", err)
		return err
	}
	yamlFile = "k8s-cluster-autoscaler.yml"
	yaml, err = template.New(yamlFile).ParseFiles(yamlFile)
	if err != nil {
		logrus.Error("Error parsing YAML file:", err)
		return err
	}

	f, err = os.Create("/var/lib/rancher/rke2/server/manifests/k8s-cluster-autoscaler.yml")
	if err != nil {
		logrus.Error("Error creating k8s-cluster-autoscaler.yml file:", err)
		return err
	}
	defer f.Close()

	err = yaml.Execute(f, cluster)
	if err != nil {
		logrus.Error("Error executing YAML template:", err)
		return err
	}

	//Cinder Storage
	yamlFile = "k8s-cinder-for-storage.yml"
	yaml, err = template.New(yamlFile).ParseFiles(yamlFile)
	if err != nil {
		logrus.Error("Error parsing YAML file:", err)
		return err
	}

	f, err = os.Create("/var/lib/rancher/rke2/server/manifests/k8s-cinder-for-storage.yml")
	if err != nil {
		logrus.Error("Error creating k8s-cinder-for-storage.yml file:", err)
		return err
	}
	defer f.Close()

	err = yaml.Execute(f, cluster)
	if err != nil {
		logrus.Error("Error executing YAML template:", err)
		return err
	}

	yamlFile = "k8s-vke-cluster-agent.yml"
	yaml, err = template.New(yamlFile).ParseFiles(yamlFile)
	if err != nil {
		logrus.Error("Error parsing YAML file:", err)
		return err
	}

	f, err = os.Create("/var/lib/rancher/rke2/server/manifests/k8s-vke-cluster-agent.yml")
	if err != nil {
		logrus.Error("Error creating k8s-vke-cluster-agent.yml file:", err)
		return err
	}

	err = yaml.Execute(f, cluster)
	if err != nil {
		logrus.Error("Error executing YAML template:", err)
		return err
	}

	err = EditCanalHelmChartConfig()
	if err != nil {
		logrus.Error("Error editing Canal Helm chart config:", err)
		return err
	}

	return nil
}
