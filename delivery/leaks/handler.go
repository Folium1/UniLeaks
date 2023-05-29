package leaks

import (
	"html/template"
	leaksRepository "uniLeaks/leaks/repository"
	leaksService "uniLeaks/leaks/service"
)

type LeaksHandler struct {
	tmpl        *template.Template
	leakService *leaksService.Service
}

func New(tmpl *template.Template) *LeaksHandler {
	leakServ := leaksService.New(leaksRepository.New())
	return &LeaksHandler{tmpl: tmpl, leakService: leakServ}
}
