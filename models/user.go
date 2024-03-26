package models

type User struct {
	Name string   `yaml:"name"`
	User UserData `yaml:"user"`
}

type UserData struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}
