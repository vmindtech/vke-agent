package models

type HelmChartConfig struct {
	APIVersion string              `yaml:"apiVersion"`
	Kind       string              `yaml:"kind"`
	Metadata   HelmChartMetadata   `yaml:"metadata"`
	Spec       HelmChartConfigSpec `yaml:"spec"`
}

type HelmChartMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type HelmChartConfigSpec struct {
	TargetNamespace string                 `yaml:"targetNamespace"`
	Bootstrap       bool                   `yaml:"bootstrap"`
	Set             map[string]interface{} `yaml:"set"`
}
