package main

import (
	"flag"
	"fmt"
)

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

	if err := utils.updateSystem(); err != nil {
		fmt.Println("System update error:", err)
		return
	}

	if err := utils.createDirectory("/etc/rancher/rke2"); err != nil {
		fmt.Println("Indexing error:", err)
		return
	}
	if err := utils.rke2Config(*initialize, *serverAddress, *rke2AgentType, *rke2Token, *tlsSan); err != nil {
		fmt.Println("Config creation error:", err)
		return
	}

	if err := utils.rke2Install(*kubeversion, *rke2AgentType); err != nil {
		fmt.Println("RKE2 installation error:", err)
		return
	}

	if err := utils.rke2ServiceEnable(*rke2AgentType); err != nil {
		fmt.Println("Service enabled error:", err)
		return
	}
	if err := utils.rke2ServiceStart(*rke2AgentType); err != nil {
		fmt.Println("Service initialization error:", err)
		return
	}
	if err := utils.pushRKE2Config(*initialize, *rke2AgentType, *serverAddress, *rke2ClusterName, *rke2ClusterUUID, *rke2AgentVKEAPIEndpoint, *rke2AgentVKEAPIAuthToken); err != nil {
		fmt.Println("Pushing RKE2 config error:", err)
		return
	}

	fmt.Println("Process completed.")
}
