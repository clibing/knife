package email

import (
	"crypto/tls"
	"fmt"

	"log"
	"net"
	"net/smtp"
	"testing"
	"time"

	gomail "gopkg.in/gomail.v2"
)

func Test163com(t *testing.T) {
	sender := "wmsjhappy@163.com"

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", sender, "clibing")

	// msg.SetHeader("From", sender)
	// msg.SetHeader("To", msg.FormatAddress("458914534@qq.com", "458914534"))
	msg.SetHeader("To", "458914534@qq.com")
	msg.SetHeader("Subject", "测试邮件")
	msg.SetBody("text/html", "<b>This is the body of the mail</b>")
	// msg.Attach("/home/User/cat.jpg")

	n := gomail.NewDialer("smtp.163.com", 465, sender, "设置设备密码")

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}
}

func TestEmail163(t *testing.T) {
	smtp_host := "smtp.163.com"
	smtp_port := 465
	senders_email := "wmsjhappy@163.com"
	senders_password := "设置设备密码"

	recipient_email := "wmsjhappy@126.com"

	type Pair struct {
		Key   string
		Value string
	}

	header := make([]*Pair, 0)
	header = append(header,
		&Pair{
			Key:   "From",
			Value: "<" + senders_email + ">",
		},
		&Pair{
			Key:   "Sender",
			Value: "knife",
		},
		&Pair{
			Key:   "To",
			Value: recipient_email,
		},
		&Pair{
			Key:   "Subject",
			Value: "邮件标题",
		},
		&Pair{
			Key:   "Content-Type",
			Value: "text/html; charset=UTF-8",
		},
	)

	body := "<html><body><div style='color: red; font-size: 12px;'>我是一封测试电子邮件! " + time.Now().Format("2006-01-02 15:04:05") + "</div></body></html>"

	message := ""
	for _, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", v.Key, v.Value)
	}

	message += "\r\n" + body

	auth := smtp.PlainAuth("", senders_email, senders_password, smtp_host)

	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", smtp_host, smtp_port),
		auth,
		senders_email,
		[]string{recipient_email},
		[]byte(message),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent successfully!")

}

func TestEmail(t *testing.T) {
	smtp_host := "smtp.mxhichina.com"
	smtp_port := 465
	senders_email := "service@linuxcrypt.cn"
	senders_password := "账号对应的密码"

	recipient_email := "wmsjhappy@163.com"

	// message := []byte("To: " + recipient_email + "\r\n" +
	// 	"Subject: Go SMTP Test\r\n" +
	// 	"\r\n" +
	// 	"Hello,\r\n" +
	// 	"This is a test email sent from Go!\r\n")

	type Pair struct {
		Key   string
		Value string
	}

	header := make([]*Pair, 0)
	header = append(header,
		&Pair{
			Key:   "From",
			Value: "<" + senders_email + ">",
		},
		&Pair{
			Key:   "Sender",
			Value: "knife",
		},
		&Pair{
			Key:   "To",
			Value: recipient_email,
		},
		&Pair{
			Key:   "Subject",
			Value: "邮件标题",
		},
		&Pair{
			Key:   "Content-Type",
			Value: "text/html; charset=UTF-8",
		},
	)

	body := "<html><body><div style='color: red; font-size: 12px;'>我是一封测试电子邮件! " + time.Now().Format("2006-01-02 15:04:05") + "</div></body></html>"

	message := ""
	for _, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", v.Key, v.Value)
	}

	message += "\r\n" + body

	auth := smtp.PlainAuth("", senders_email, senders_password, smtp_host)

	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", smtp_host, smtp_port),
		auth,
		senders_email,
		[]string{recipient_email},
		[]byte(message),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent successfully!")

}

// return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 参考net/smtp的func SendMail()
// 使用net.Dial连接tls（SSL）端口时，smtp.NewClient()会卡住且不提示err
// len(to)>1时，to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
