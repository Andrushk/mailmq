package context

import (
	"io/ioutil"
	"encoding/json"
)

type AppConfig struct {
	// true - бесшумный режим, сообщения из очереди записывать с лог, по почте не отправлять
	SilentMode bool `json:"silent_mode"`

	// настройки очереди
	RabbitMQHost string `json:"rabbit_mq_host"`
	RabbitMQQueueName string `json:"rabbit_mq_queue_name"`

	// настройки почтового сервера
	MailHost string `json:"mail_host"`
	MailUserName string `json:"mail_user_name"`
	MailPassword string `json:"mail_password"`
}

func LoadConfig(path string) (*AppConfig, error)  {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &AppConfig{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}