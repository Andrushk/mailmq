package bl

import (
	"github.com/andrushk/mailmq/context"
	"github.com/andrushk/mailmq/dto"
	"fmt"
)

// Реализация интерфейса Sender для тестирования
// отправляет сообщения в лог
type SendToLog struct {
	Log context.Logger
}

func (m *SendToLog) Send(task dto.Task) error {
	m.Log.Info(fmt.Sprintf("\r\nFrom: %s\r\nTo: %s\r\nSubject: %s\r\nMessage: %s\r\n",
		task.From,
		task.To,
		task.Subject,
		task.Message))
	return nil
}
