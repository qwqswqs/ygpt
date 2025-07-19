package sms

import (
	"encoding/json"
	"errors"
	"github.com/volcengine/volc-sdk-golang/service/sms"
	"log"
	"ygpt/server/global"
)

// sendSmsTemplateParam 忘记密码
type sendSmsTemplateParam struct {
	Code string `json:"code"`
}

// sendSmsTemplateParam2 忘记密码
type sendSmsTemplateParam2 struct {
	Message string `json:"message"`
}

// PasswordTpl 忘记密码
type PasswordTpl struct {
	Code string `json:"code"`
}

// AuthFail 企业用户认证认证失败
type AuthFail struct {
	CompanyName []string `json:"companyType"`
}

// AuthFail2 企业用户认证认证失败
type AuthFail2 struct {
}

// Subscribe 推送订阅通知
type Subscribe struct {
	//订阅用户名称
	User string `json:"user"`
	//平台推荐的信息名称
	Name string `json:"name"`
}

// ForgetPasswordTpl 忘记密码
func ForgetPasswordTpl(userId int, code string, phone string) error {
	dataMap := make(map[string]string, 1)
	dataMap["code"] = code

	cj, _ := json.Marshal(dataMap)
	return SendSmsCommon(global.GVA_CONFIG.Sms.TemplateId, string(cj), phone)
}

// SendSmsCommon 短信发送通用方法
// templateId 短信模板id
// templateParam 短信参数json字符串
// phone 手机号
func SendSmsCommon(templateId, templateParam, phone string) error {
	sms.DefaultInstance.Client.SetAccessKey(global.GVA_CONFIG.Sms.AccessKey)
	sms.DefaultInstance.Client.SetSecretKey(global.GVA_CONFIG.Sms.SecretKey)
	req := &sms.SmsRequest{
		SmsAccount:    global.GVA_CONFIG.Sms.SmsAccount,
		Sign:          global.GVA_CONFIG.Sms.SmsSign,
		TemplateID:    templateId,
		TemplateParam: templateParam,
		PhoneNumbers:  phone,
		Tag:           global.GVA_CONFIG.Sms.Tag,
	}
	result, _, err := sms.DefaultInstance.Send(req)
	if err != nil {
		log.Println(err)
		return err
	} else {
		if result.ResponseMetadata.Error != nil {
			return errors.New(result.ResponseMetadata.Error.Message)
		}
		return nil
	}
}

// SubscribeTpl  SubscribeTpl
func SubscribeTpl(phone, user, name string) error {
	dataMap := make(map[string]string, 1)
	dataMap["name"] = name
	dataMap["user"] = user
	cj, _ := json.Marshal(dataMap)
	return SendSmsCommon(global.GVA_CONFIG.Sms.TemplateSubId, string(cj), phone)
}

//sms.ForgetPasswordTpl(1, "123456", "19508199151")
//sms.AuthFailTpl("19508199151")
//sms.AuthSuccessTpl("19508199151")
//sms.RegTpl("19508199151")
//sms.SubscribeTpl("19508199151", "gzw", "nihao ")
