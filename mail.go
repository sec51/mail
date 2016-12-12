package mail

import (
	"crypto/tls"
	"errors"

	"github.com/sec51/goconf"

	gomail "gopkg.in/gomail.v2"
)

type EmailConfig struct {
	// The hostname or ip address of the mail server
	Server string

	// The mail server submission port
	Port int

	// the From email: can be in the form: "Info <info@sec51.com>"
	From string

	// The username to authenticate
	Username string

	// The password to authenticate
	Password string
}

func (c *EmailConfig) Send(body, to, subject string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", c.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpServer := goconf.AppConf.String("smtp::server")
	if smtpServer == "" {
		return errors.New("SMTP server not configured,set: smtp.server = HOSTNAME in the config file")
	}

	smtpUser := goconf.AppConf.DefaultString("smtp::username", c.Username)
	if smtpUser == "" {
		return errors.New("SMTP username not configured, set: smtp::username = NAME in the config file")
	}

	smtpPass := goconf.AppConf.DefaultString("smtp::password", c.Password)
	if smtpPass == "" {
		return errors.New("SMTP password not configured, set: smtp::password = PASSWORD in the config file")
	}

	smtpPort := goconf.AppConf.DefaultInt("smtp::port", c.Port)

	d := gomail.NewPlainDialer(smtpServer, smtpPort, smtpUser, smtpPass)
	d.SSL = true
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpServer,
	}
	d.TLSConfig = tlsconfig

	// Send the email to Bob, Cora and Dan
	return d.DialAndSend(m)
}
