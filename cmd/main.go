package main

import (
	"leaks/pkg/config"
	delivery "leaks/pkg/http"
)

func main() {
	config.InitMYSQL()
	handler := delivery.New()
	handler.StartServer()
}
