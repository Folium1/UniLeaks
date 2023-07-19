package http

import (
	leaksService "leaks/pkg/leaks/service"
	logg "leaks/pkg/logger"
)

var logger = logg.NewLogger()

type LeaksHandler struct {
	leakService *leaksService.Service
}

func New() *LeaksHandler {
	leakServ := leaksService.New()
	return &LeaksHandler{leakService: leakServ}
}
