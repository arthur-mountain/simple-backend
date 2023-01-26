package smtp

import (
	"net/smtp"
)

type sender struct {
	email    string
	password string
}

type smtpConfig struct {
	host      string
	port      string
	sender    sender
	receivers []string
	auth      smtp.Auth
}

func Init(host string, port string, sender sender, receivers []string) *smtpConfig {
	auth := smtp.PlainAuth("", sender.email, sender.password, host)

	return &smtpConfig{
		host,
		port,
		sender,
		receivers,
		auth,
	}
}

func (s *smtpConfig) Update(key string, value interface{}) *smtpConfig {
	switch key {
	case "host":
		s.host = value.(string)
	case "port":
		s.port = value.(string)
	case "sender":
		s.sender = value.(sender)
	case "receivers":
		s.receivers = value.([]string)
	}

	return s
}

func (s *smtpConfig) SendEmail(message []byte) error {
	err := smtp.SendMail(s.host+":"+s.port, s.auth, s.sender.email, s.receivers, message)

	if err != nil {
		return err
	}

	return nil
}

// func (s *smtpConfig) SendEmailWithAttach(message []byte) error {
// 	err := smtp.SendMail(s.host+":"+s.port, s.auth, s.sender.email, s.receivers, message)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
