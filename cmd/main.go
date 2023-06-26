package main

import (
	"leaks/config"
	"leaks/delivery"
)

func main() {
	config.InitMYSQL()
	handler := delivery.New()
	handler.StartServer()
}
