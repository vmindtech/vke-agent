package main

type InitMaster struct {
	NodeName      string
	Token         string
	TlsSan        string
	Initialize    bool
	ServerAddress string
	Rke2AgentType string
}

type KubeConfig struct {
	APIVersion     string    `yaml:"apiVersion"`
	Clusters       []Cluster `yaml:"clusters"`
	Contexts       []Context `yaml:"contexts"`
	CurrentContext string    `yaml:"current-context"`
	Kind           string    `yaml:"kind"`
	Preferences    struct{}  `yaml:"preferences"`
	Users          []User    `yaml:"users"`
}

type Cluster struct {
	Cluster ClusterData `yaml:"cluster"`
	Name    string      `yaml:"name"`
}

type ClusterData struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type Context struct {
	Context ContextData `yaml:"context"`
	Name    string      `yaml:"name"`
}

type ContextData struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type User struct {
	Name string   `yaml:"name"`
	User UserData `yaml:"user"`
}

type UserData struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}
