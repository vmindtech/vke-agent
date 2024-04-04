package models

type Config struct {
	ServerAddress            string
	Kubeversion              string
	TLSSan                   string
	Initialize               bool
	RKE2Token                string
	RKE2AgentType            string
	RKE2NodeLabel            string
	RKE2ClusterName          string
	RKE2ClusterUUID          string
	RKE2ClusterProjectUUID   string
	RKE2AgentVKEAPIEndpoint  string
	RKE2AgentVKEAPIAuthToken string
	VkeCloudAuthURL          string
	ClusterAutoscalerVersion string
	CloudProviderVkeVersion  string
	ApplicationCredentialID  string
	ApplicationCredentialKey string
}
