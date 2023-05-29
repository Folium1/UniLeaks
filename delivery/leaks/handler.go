package leaks

import (
	"html/template"
	leaksRepository "uniLeaks/leaks/repository"
	leaksService "uniLeaks/leaks/service"

	"github.com/gin-gonic/gin"
)

type LeaksHandler struct {
	router      *gin.Engine
	tmpl        *template.Template
	leakService *leaksService.Service
}

func New(router *gin.Engine, tmpl *template.Template) *LeaksHandler {
	leakServ := leaksService.New(leaksRepository.New())
	return &LeaksHandler{router: router, tmpl: tmpl, leakService: leakServ}
}
