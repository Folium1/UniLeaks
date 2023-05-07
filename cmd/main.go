package main

import (
	"uniLeaks/config"
	delivery "uniLeaks/delivery/http"
)

func main() {
	config.InitMYSQL()
	handler := delivery.New()
	handler.StartServer()
}
