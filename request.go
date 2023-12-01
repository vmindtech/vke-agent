package main

type SendKubeConfigRequest struct {
	ClusterID  string `json:"clusterId"`
	KubeConfig string `json:"kubeconfig"`
}
