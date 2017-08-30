package bl

import (
	"net/smtp"
	"github.com/andrushk/mailmq/context"
	"fmt"
	"github.com/andrushk/mailmq/dto"
)

// Реализация интерфейса Sender
// отправляет сообщения по электронной почте
type MailSender struct {
	ctx *context.AppContext
	auth smtp.Auth
}

func CreateMailSender(ctx *context.AppContext) *MailSender {
	// Set up authentication information.
	auth := smtp.PlainAuth("",
		ctx.Cgf.MailUserName,
		ctx.Cgf.MailPassword,
		ctx.Cgf.MailHost)
	return &MailSender{ctx:ctx, auth:auth}
}

func (ms *MailSender) Send(task dto.Task) error {
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", task.Subject, task.Message))
	return smtp.SendMail(ms.ctx.Cgf.MailHost, ms.auth, task.From, task.To, msg)
}
