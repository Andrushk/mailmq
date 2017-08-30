package bl

import "github.com/andrushk/mailmq/dto"

type Sender interface {
	Send(task dto.Task) error
}
