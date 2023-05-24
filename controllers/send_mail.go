package controllers

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func SendMail(RandomCode, Receiver string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "testing.go.mail.sender@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", Receiver)

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	// Set E-Mail body. You can set plain text or html with text/html
	var bodyMessage string = "Enter the code : " + RandomCode
	m.SetBody("text/plain", bodyMessage)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "testing.go.mail.sender@gmail.com", "xjgpwjsxceeutrhx")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		// panic(err)
		return
	}

	return
}
