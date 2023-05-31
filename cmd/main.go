package main

import (
	"uniLeaks/config"
	delivery "uniLeaks/delivery"
)

func main() {
	config.InitMYSQL()
	handler := delivery.New()
	handler.StartServer()
}
