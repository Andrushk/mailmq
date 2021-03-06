package bl

import (
	"net/smtp"
	"github.com/andrushk/mailmq/context"
	"fmt"
	"github.com/andrushk/mailmq/dto"
	"crypto/tls"
	"strings"
	"net/mail"
	"github.com/andrushk/mailmq/consts"
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
	ctx.Log.Info(fmt.Sprintf(consts.EmailServerNameConfirmation, smtpServer.ServerName()))

	return &MailSender{
		sender:     ctx.Cgf.MailUserName,
		auth:       smtp.PlainAuth("", ctx.Cgf.MailUserName, ctx.Cgf.MailPassword, smtpServer.Host),
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

	// Есть замечательный и простой способ отправить мэйл:
	// smtp.SendMail(ms.smtpServer.ServerName(), ms.auth, ms.sender, task.To, []byte(msg.BuildMessage()))
	// но на ubuntu им не получается воспользоваться, возникает ошибка: x509: certificate signed by unknown authority
	// чтобы ее забороть необходимо для TlsConfig установить параметр: InsecureSkipVerify: true
	// поэтому приходится городить огород ниже:

	conn, err := smtp.Dial(ms.smtpServer.ServerName())
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.StartTLS(ms.smtpServer.TlsConfig)
	if err != nil {
		return err
	}

	// Auth
	if err = conn.Auth(ms.auth); err != nil {
		return err
	}

	// To && From
	if err = conn.Mail(ms.sender); err != nil {
		return err
	}

	for _, addr := range task.To {
		if err = conn.Rcpt(addr); err != nil {
			return err
		}
	}

	// Data
	w, err := conn.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg.BuildMessage()))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return conn.Quit()
}
