package config

type CloudpodsConfig struct {
	AuthURL      string `yaml:"authUrl" mapstructure:"authUrl" json:"authUrl"`
	JumpUrl      string `yaml:"jumpUrl" mapstructure:"jumpUrl" json:"jumpUrl"`
	Timeout      int    `yaml:"timeout" mapstructure:"timeout" json:"timeout"`
	Debug        bool   `yaml:"debug" mapstructure:"debug" json:"debug"`
	Insecure     bool   `yaml:"insecure" mapstructure:"insecure" json:"insecure"`
	Uname        string `yaml:"uname" mapstructure:"uname" json:"uname"`
	Passwd       string `yaml:"passwd" mapstructure:"passwd" json:"passwd"`
	DomainName   string `yaml:"domainName" mapstructure:"domainName" json:"domainName"`
	TenantName   string `yaml:"tenantName" mapstructure:"tenantName" json:"tenantName"`
	TenantDomain string `yaml:"tanantDomain" mapstructure:"tanantDomain" json:"tanantDomain"`
	Region       string `yaml:"region" mapstructure:"region" json:"region"`
	EndpointType string `yaml:"endpointType" mapstructure:"endpointType" json:"endpointType"`
}
