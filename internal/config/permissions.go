package config

type Permissions struct {
	Index  *bool `yaml:"index"`
	View   bool  `yaml:"view"`
	Deploy bool  `yaml:"deploy"`
}
