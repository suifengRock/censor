package main

import (
	"censor/service"
)

func block() {
	block := make(chan bool)
	<-block
}

func main() {
	service.PushRequestUrl("")

	block()
}
