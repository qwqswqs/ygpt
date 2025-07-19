package email

import (
	"github.com/go-gomail/gomail"
	"ygpt/server/global"
)

// SendEmail 邮箱发送
func SendEmail(from string, to string, body string, host string, port int, userName string, password string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to, to)
	m.SetAddressHeader("Cc", from, "admin")
	m.SetHeader("Subject", "浮点运算")
	m.SetBody("text/html", body)
	d := gomail.NewDialer(host, port, userName, password)
	return d.DialAndSend(m)
}

//"admin@slquan.com", user.Email, message, "smtp.feishu.cn", 587, "admin@slquan.com", "1rfMNeCrjVcvla2f"

// SendEmail2 邮箱发送
func SendEmail2(to string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", global.GVA_CONFIG.Email.From)
	m.SetHeader("To", to, to)
	m.SetAddressHeader(global.GVA_CONFIG.Email.From, global.GVA_CONFIG.Email.From, "浮点运算运营商")
	m.SetHeader("Subject", "浮点运算")
	m.SetBody("text/html", body)
	d := gomail.NewDialer(global.GVA_CONFIG.Email.Host, global.GVA_CONFIG.Email.Port, global.GVA_CONFIG.Email.From, global.GVA_CONFIG.Email.Secret)
	return d.DialAndSend(m)
}
