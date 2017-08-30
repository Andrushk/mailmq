package main

import (
	"fmt"
	"log"
)

// Реализация context/Logger
type AppLog struct {
}

func (l AppLog) Info(comment interface{}) {
	log.Printf("%s", comment)
}

func (l AppLog) Fatal(comment interface{}, err error) string {
	mess := fmt.Sprintf("%s, original error: %s", comment, err)
	log.Fatalf(mess)
	return mess
}

func (l AppLog) Warning(comment interface{}, err error) {
	log.Printf("%s, original error: %s", comment, err)
}
