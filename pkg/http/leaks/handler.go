package leaks

import (
	"html/template"
	leaksRepository "leaks/pkg/leaks/repository"
	leaksService "leaks/pkg/leaks/service"
	"leaks/pkg/logger"
)

var logg = logger.NewLogger()

type LeaksHandler struct {
	tmpl        *template.Template
	leakService *leaksService.Service
}

func New(tmpl *template.Template) *LeaksHandler {
	leakServ := leaksService.New(leaksRepository.New())
	return &LeaksHandler{tmpl: tmpl, leakService: leakServ}
}
