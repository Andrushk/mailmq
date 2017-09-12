# Отслает по email сообщения из очереди RabbitMQ

Настройка (что читать и куда посылать) производится через конфиг-файл. При запуске 'mailmq' необходимо указать параметром имя этого файла. Для примера в проект включен 'mailmq.cfg.example'. Пример его использования: 'mailmq mailmq.cfg.example'.

Описание 'mailmq.cfg.example':

```
{
  // true - отправки по email не происходит, сообщение из очереди читается и выводится в лог
  // false - нормальный рабочий режим, сообщения отправляются по email
  "silent_mode": false,
  
  // имя сервера с RabbitMQ, в примере случай, когда mailmq и RabbitMQ запущены на одной машине
  "rabbit_mq_host": "amqp://",
  
  // имя очереди, должно совпадать с источником сообщений
  "rabbit_mq_queue_name": "MailsQueue",
  
  // далее настройки почтового севера
  "mail_host" : "smtp.mail.ru",
  "mail_port" : "25",
  "mail_user_name": "mail@mail.ru",
  "mail_password": "12345"
}
```
