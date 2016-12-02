package main

import (
	"bytes"
	"log"
	"net/smtp"
	"strings"
	"text/template"
)

const (
	// email template
	message = `To: {{.M.ToString}}
From: Mikrotik Firmware<{{.M.From}}>
Subject: New Mikrotik Firmware {{.E.Version}} version is available

{{.E.Info}}`
)

// mail provides logic to send email notifications.
type mail struct {
	To       []string // mail 'To' header
	From     string   // mail 'From' header
	Addr     string   // SMTP address (example: smtp.gmail.com:25)
	Username string   // SMTP Auth username
	Password string   // SMTP Auth password

	auth smtp.Auth
}

// newMail intializes new mail object.
func newMail(to, from, addr, username, password string) *mail {
	if len(to) == 0 {
		log.Fatal("The 'to' flag is required.")
	}
	if len(from) == 0 {
		log.Fatal("The 'from' flag is required.")
	}
	if len(addr) == 0 {
		log.Fatal("The 'addr' flag is required.")
	}
	if len(username) == 0 {
		username = from
	}
	if len(password) == 0 {
		log.Fatal("The 'password' flag is required.")
	}
	m := &mail{
		To:       strings.Split(to, ","),
		From:     from,
		Addr:     addr,
		Username: username,
		Password: password,
	}
	m.auth = smtp.PlainAuth("", m.Username, m.Password, m.hostname())
	return m
}

// hostname gets hostname for the address string.
func (m mail) hostname() string {
	args := strings.Split(m.Addr, ":")
	return args[0]
}

// ToString converts To slice to string.
func (m mail) ToString() string {
	return strings.Join(m.To, ",")
}

// Notify sends notification about new Mikrotik firmware version available.
func (m mail) Notify(ver, info string) {
	data := struct {
		M mail
		E map[string]string
	}{
		M: m,
		E: map[string]string{"Version": ver, "Info": info},
	}
	var bs bytes.Buffer
	err := template.Must(template.New("").Parse(message)).Execute(&bs, data)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(bs.String(), "\n")
	body := strings.Join(lines, "\r\n") + "\r\n"
	err = smtp.SendMail(m.Addr, m.auth, m.From, m.To, []byte(body))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Notifications are sent to %s", m.ToString())
	}
}
