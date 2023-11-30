package main

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
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
	Name                     string `yaml:"name"`
}

type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
	Name    string `yaml:"name"`
}

type User struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
	Name                  string `yaml:"name"`
}
