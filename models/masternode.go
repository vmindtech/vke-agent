package models

type InitMaster struct {
	NodeName                    string
	Token                       string
	TlsSan                      string
	Initialize                  bool
	ServerAddress               string
	Rke2AgentType               string
	Rke2NodeLabel               []string
	VmindCloudAuthURL           string
	ApplicationCredentialID     string
	ApplicationCredentialSecret string
	ClusterAutoscalerVersion    string
	CloudProviderVkeVersion     string
}
