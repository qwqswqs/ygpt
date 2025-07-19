package config

type SLMinioConfig struct {
	EndPoint   string `yaml:"EndPoint" mapstructure:"EndPoint" json:"EndPoint"`
	AccessKey  string `yaml:"AccessKey" mapstructure:"AccessKey" json:"AccessKey"`
	SecretKey  string `yaml:"SecretKey" mapstructure:"SecretKey" json:"SecretKey"`
	BucketName string `yaml:"BucketName" mapstructure:"BucketName" json:"BucketName"`
	UseSSL     int    `yaml:"UseSSL" mapstructure:"UseSSL" json:"UseSSL"`
}
