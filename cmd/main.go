package main

import (
	leaksApi "leaks/pkg/api"
	"leaks/pkg/config"
)

func main() {
	config.InitMysqlTables()
	api := leaksApi.NewApiHandler()
	api.StartServer()
}
