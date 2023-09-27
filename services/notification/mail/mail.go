package mail

import (

	//"gitlab.com/u-v/dash/server/services/assert"

	"crypto/tls"

	"github.com/go-gomail/gomail"
	"gitlab.com/mefit/mefit-server/services/notification"
	"gitlab.com/mefit/mefit-server/utils/assert"
	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/initializer"
	"gitlab.com/mefit/mefit-server/utils/log"
)

var (
	// dialer *gomail.Dialer

	host       string
	port       int
	user, pass string
)

type manager struct{}

func Send(notif *notification.Notification) {
	if !notif.Validate() {
		log.Logger().Panic("mail is not valid")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", notif.From)
	m.SetHeader("To", notif.To...)
	m.SetHeader("Subject", notif.Subject)
	m.SetBody("text/html", notif.Body)

	err := dialer.DialAndSend(m)
	assert.Nil(err)
}

var dialer *gomail.Dialer

func init() {
	initializer.Register(manager{})
}

func (manager) Initialize() func() {
	host = config.Config().GetDefaultString("smtp_host", "smtp.gmail.com")
	user = config.Config().GetDefaultString("smtp_username", "")
	pass = config.Config().GetDefaultString("smtp_password", "")
	port = config.Config().GetDefaultInt("smtp_port", 465)
	println("user", user)
	println("pass", pass)
	dialer = gomail.NewPlainDialer(host, port, user, pass)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false}

	// _, err := dialer.Dial()
	// assert.Nil(err)

	return nil
}
