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

	// whether to make a TLS connection
	UseTLS bool

	// whether to use username and password for authentication
	UseAUTH bool
}

func (c *EmailConfig) Send(body, to, subject string) error {

	m := gomail.NewMessage()
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	from := c.From
	if from == "" {
		from = goconf.AppConf.String("smtp::from")
		if from == "" {
			return errors.New("From email address not configured - set: smtp::from = FROM in the config file")
		}
	}

	m.SetHeader("From", from)

	smtpServer := c.Server
	if smtpServer == "" {
		smtpServer = goconf.AppConf.String("smtp::server")
		if smtpServer == "" {
			return errors.New("SMTP server not configured - set: smtp::server = HOSTNAME in the config file")
		}
	}

	smtpUser := c.Username
	if smtpUser == "" {
		smtpUser = goconf.AppConf.DefaultString("smtp::username", c.Username)
		if smtpUser == "" {
			//return errors.New("SMTP username not configured, set: smtp::username = NAME in the config file")
		}
	}

	smtpPass := c.Password
	if smtpPass == "" {
		smtpPass = goconf.AppConf.DefaultString("smtp::password", c.Password)
		if smtpPass == "" {
			//return errors.New("SMTP password not configured, set: smtp::password = PASSWORD in the config file")
		}
	}

	smtpPort := c.Port
	if smtpPort == 0 {
		smtpPort = goconf.AppConf.DefaultInt("smtp::port", c.Port)
		if smtpPort == 0 {
			return errors.New("SMTP port not configured, set: smtp::port = PORT in the config file")
		}
	}

	if c.UseAUTH {
		d := gomail.NewPlainDialer(smtpServer, smtpPort, smtpUser, smtpPass)
		d.SSL = c.UseTLS

		if d.SSL {
			tlsconfig := &tls.Config{
				InsecureSkipVerify: false,
				ServerName:         smtpServer,
			}
			d.TLSConfig = tlsconfig
		}

		// Send the email to Bob, Cora and Dan
		return d.DialAndSend(m)
	}

	d := gomail.Dialer{Host: smtpServer, Port: smtpPort}
	return d.DialAndSend(m)

}
