package bl

import (
	"net/smtp"
	"github.com/andrushk/mailmq/context"
	"fmt"
	"github.com/andrushk/mailmq/dto"
	"crypto/tls"
	"strings"
	"net/mail"
)

// Реализация интерфейса Sender
// отправляет сообщения по электронной почте
type MailSender struct {
	sender     string
	auth       smtp.Auth
	smtpServer SmtpServer
}

// Настройки одного письма
type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

// Настройки сервера
type SmtpServer struct {
	Host      string
	Port      string
	TlsConfig *tls.Config
}

// Получить экземпляр отправлялки писем
func CreateMailSender(ctx *context.AppContext) *MailSender {
	smtpServer := SmtpServer{Host: ctx.Cgf.MailHost, Port: ctx.Cgf.MailPort}
	smtpServer.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.Host,
	}

	auth := smtp.PlainAuth("", ctx.Cgf.MailUserName, ctx.Cgf.MailPassword, smtpServer.Host)

	return &MailSender{
		sender: ctx.Cgf.MailUserName,
		auth: auth,
		smtpServer: smtpServer}
}

func (s *SmtpServer) ServerName() string {
	return s.Host + ":" + s.Port
}

func (em *Mail) BuildMessage() string {
	header := ""
	header += fmt.Sprintf("From: %s\r\n", em.Sender)
	if len(em.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(em.To, ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", em.Subject)
	header += "\r\n" + em.Body

	return header
}

func (em *Mail) SenderAddress() string  {
	return (mail.Address{"", em.Sender}).Address
}

func (ms *MailSender) Send(task dto.Task) error {
	msg := Mail{}
	msg.Sender = ms.sender
	msg.To = task.To
	msg.Subject = task.Subject
	msg.Body = task.Message

	return smtp.SendMail(ms.smtpServer.ServerName(), ms.auth, ms.sender, task.To, []byte(msg.BuildMessage()))
}
