package config

type Sms struct {
	AccessKey              string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`
	SecretKey              string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	SmsAccount             string `mapstructure:"sms-account" json:"sms-account" yaml:"sms-account"`
	TemplateId             string `mapstructure:"template-id" json:"template-id" yaml:"template-id"`
	TemplateSubId          string `mapstructure:"template-sub-id" json:"template-sub-id"  yaml:"template-sub-id"`
	TemplateAuthFailId     string `mapstructure:"template-auth-fail-id" json:"template-auth-fail-id"  yaml:"template-auth-fail-id"`
	TemplateAuthSuccessId  string `mapstructure:"template-auth-success-id" json:"template-auth-success-id"  yaml:"template-auth-success-id"`
	TemplateRegisSuccessId string `mapstructure:"template-regis-success-id" json:"template-regis-success-id"  yaml:"template-regis-success-id"`
	TemplateAlertId        string `mapstructure:"template-alert-id" json:"template-alert-id"  yaml:"template-alert-id"`
	SmsSign                string `mapstructure:"sms-sign" json:"sms-sign"  yaml:"sms-sign"`
	Tag                    string `mapstructure:"tag" json:"tag" yaml:"tag"`
}
