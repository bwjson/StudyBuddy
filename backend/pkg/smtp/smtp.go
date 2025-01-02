package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPServer struct {
	Host     string
	Port     string
	From     string
	Password string
}

func NewSMTPServer(host string, port string, from string, password string) *SMTPServer {
	return &SMTPServer{host, port, from, password}
}

func (server *SMTPServer) SendEmail(email string, subject string, message string) error {
	auth := smtp.PlainAuth("", server.From, server.Password, server.Host)

	htmlMessage := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>%s</title>
		</head>
		<body>
			<h1>%s</h1>
			<p>%s</p>
		</body>
		</html>
	`, subject, subject, message)

	headers := map[string]string{
		"From":         server.From,
		"To":           email,
		"Subject":      subject,
		"Content-Type": "text/html; charset=UTF-8",
	}

	var emailMessage []string
	for key, value := range headers {
		emailMessage = append(emailMessage, fmt.Sprintf("%s: %s", key, value))
	}
	emailMessage = append(emailMessage, "")

	fullMessage := strings.Join(emailMessage, "\r\n") + htmlMessage

	err := smtp.SendMail(server.Host+":"+server.Port, auth, server.From, []string{email}, []byte(fullMessage))
	if err != nil {
		return err
	}

	return nil
}

func (server *SMTPServer) SendVerifyingEmail(email string, subject string, url string) error {
	auth := smtp.PlainAuth("", server.From, server.Password, server.Host)

	htmlMessage := fmt.Sprintf(`
        <html>
        <body>
            <p>Thank you for registering!</p>
            <p>Click the link below to verify your email:</p>
            <a href="%v">Link</a>
        </body>
        </html>`, url)

	headers := map[string]string{
		"From":         server.From,
		"To":           email,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var emailMessage []string
	for key, value := range headers {
		emailMessage = append(emailMessage, fmt.Sprintf("%s: %s", key, value))
	}

	emailMessage = append(emailMessage, "")
	emailMessage = append(emailMessage, htmlMessage)

	fullMessage := strings.Join(emailMessage, "\r\n")

	err := smtp.SendMail(server.Host+":"+server.Port, auth, server.From, []string{email}, []byte(fullMessage))
	if err != nil {
		return err
	}

	return nil
}
