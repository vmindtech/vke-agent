package models

type Cluster struct {
	Cluster ClusterData `yaml:"cluster"`
	Name    string      `yaml:"name"`
}

type ClusterData struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}
