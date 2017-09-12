package main

import (
	"github.com/andrushk/mailmq/bl"
	"github.com/andrushk/mailmq/consts"
	"github.com/andrushk/mailmq/context"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic(consts.AppArgumentsFailed)
	}

	appLog := &AppLog{}
	appConfig, err := context.LoadConfig(os.Args[1])
	if err != nil {
		panic(appLog.Fatal(consts.AppFailedToLoadConfig, err))
	}

	ctx := &context.AppContext{
		Log: appLog,
		Cgf: appConfig,
	}

	var sender bl.Sender
	if appConfig.SilentMode {
		// заглушка вместо отправлялки почты, пишет в лог
		sender = &bl.SendToLog{Log: appLog}
	} else {
		// отправка почтой
		sender = bl.CreateMailSender(ctx)
	}

	q := bl.CreateQueue(ctx, sender)
	q.Process()
}
