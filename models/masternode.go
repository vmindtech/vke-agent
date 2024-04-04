package models

type InitMaster struct {
	NodeName                 string
	Token                    string
	TlsSan                   string
	Initialize               bool
	ServerAddress            string
	Rke2AgentType            string
	Rke2NodeLabel            []string
	RKE2ClusterProjectUUID   string
	RKE2ClusterUUID          string
	VkeCloudAuthURL          string
	ApplicationCredentialID  string
	ApplicationCredentialKey string
	ClusterAutoscalerVersion string
	CloudProviderVkeVersion  string
}
