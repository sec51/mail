package mail

import (
	"crypto/tls"
	"errors"
	"github.com/sec51/goconf"
	gomail "gopkg.in/gomail.v2"
)

const (
	FROM_EMAIL  = "Falcon <info@sec51.com>"
	SMTP_SERVER = "mail.blackdefense.com"
	SMTP_PORT   = 587
)

func Send(body, to, subject string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", FROM_EMAIL)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpServer := goconf.AppConf.String("smtp.server")
	if smtpServer == "" {
		return errors.New("SMTP server not configured,set: smtp.server = HOSTNAME in the config file")
	}

	smtpUser := goconf.AppConf.String("smtp.user")
	if smtpUser == "" {
		return errors.New("SMTP user not configured,set: smtp.user = NAME in the config file")
	}

	smtpPass := goconf.AppConf.String("smtp.pass")
	if smtpUser == "" {
		return errors.New("SMTP password not configured,set: smtp.pass = PASSWORD in the config file")
	}

	smtpPort := goconf.AppConf.DefaultInt("smtp.port", SMTP_PORT)

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

