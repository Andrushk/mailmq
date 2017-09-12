package bl

import (
	"github.com/streadway/amqp"
	"github.com/andrushk/mailmq/context"
	"github.com/andrushk/mailmq/consts"
	"encoding/json"
	"github.com/andrushk/mailmq/dto"
	"fmt"
)

// Очередь на RabbitMQ
type TaskQueue struct {
	ctx        *context.AppContext
	sender     Sender

	// RabbitMQ
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func CreateQueue(ctx *context.AppContext, s Sender) TaskQueue {
	return TaskQueue{ctx:ctx, sender: s}
}

func (q *TaskQueue) Process() {
	q.initConnection()
	defer q.connection.Close()

	q.initChannel()
	defer q.channel.Close()

	q.initQueue()

	for m := range q.getTasks() {
		task := dto.Task{}

		if err := json.Unmarshal(m.Body, &task); err != nil {
			// ругаемся в лог и пропускаем такое сообщение
			q.ctx.Log.Warning(
				fmt.Sprintf("%s [%s]",
					consts.SenderMessageUnableToUnmarshal, m.Body), err)
		} else {
			//TODO: мб стоит тут не просто ругаться в лог и помечать сообщение как обработанное,
			//TODO: а запустить еще пару попыток отправки
			//TODO: или запилить специальную очередь для проблемных сообщений
			if err := q.sender.Send(task); err != nil {
				q.ctx.Log.Warning(consts.SenderMessageWasNotSent, err)
			} else {
				q.ctx.Log.Info(fmt.Sprintf("processed: %s", task.To))
			}
		}
		m.Ack(false)
	}
}

func (q *TaskQueue) initConnection() {
	var err error
	q.connection, err = amqp.Dial(q.ctx.Cgf.RabbitMQHost)
	q.breakIfError(consts.RabbitMQFailedToConnect, err)
}

func (q *TaskQueue) initChannel()  {
	var err error
	q.channel, err = q.connection.Channel()
	q.breakIfError(consts.RabbitMQFailedToOpenChannel, err)

	err = q.channel.Qos(
		1,     // из очереи выбираем по одной записи
		0,     // prefetch size
		false, // global
	)
	q.breakIfError(consts.RabbitMQFailedToSetQoS, err)
}

func (q *TaskQueue) initQueue() {
	var err error
	q.queue, err = q.channel.QueueDeclare(
		q.ctx.Cgf.RabbitMQQueueName, // name
		true,   // устойчивая очередь
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	q.breakIfError(consts.RabbitMQFailedToDeclareQueue, err)
}

func (q *TaskQueue) getTasks() <-chan amqp.Delivery {
	tasks, err := q.channel.Consume(
		q.queue.Name, // queue
		"", // consumer
		false, // авто-подтверждение, что задача выполнена, выкл
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
	q.breakIfError(consts.RabbitMQFailedToRegisterConsumer, err)
	return tasks
}

func (q *TaskQueue) breakIfError(comment string, err error){
	if err != nil {
		panic(q.ctx.Log.Fatal(comment, err))
	}
}