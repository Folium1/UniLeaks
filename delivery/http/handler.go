package delivery

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	authHttp "uniLeaks/auth/delivery/http"
	leak "uniLeaks/leaks"
	leakService "uniLeaks/leaks/service"
	"uniLeaks/user"
	userMysql "uniLeaks/user/repository/mysql"
	service "uniLeaks/user/service"

	"github.com/gin-gonic/gin"
)

var (
	Middleware = authHttp.New()
)

type subjects struct {
	Name  string
	Value string
}

type Handler struct {
	tmpl        *template.Template
	userService user.Repository
	leakService leak.Repository
	router      *gin.Engine
}

func New() Handler {
	tmpl := template.Must(template.ParseGlob("templates/*"))
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/templates/static", "/")
	newUserHandler := Handler{tmpl, service.New(*userMysql.New()), leakService.New(), r}
	return newUserHandler
}

func (h Handler) handleUsers() {
	userGroup := h.router.Group("/user")
	{
		userGroup.GET("/sign-up", h.getRegister)
		userGroup.GET("/login", h.getLogin)
		userGroup.POST("/sign-up", h.postRegister)
		userGroup.POST("/login", h.postLogin)
	}
}

func (h Handler) handleLeaks() {
	leaks := h.router.Group("/leaks")
	{
		leaks.GET("/", h.mainPage)
		uploadFiles := leaks.Group("/", Middleware.AuthAndRefreshMiddleware())
		{
			uploadFiles.GET("/upload-file", h.uploadFilePage)
			uploadFiles.POST("/upload-file", h.uploadFile)
		}
		getFiles := leaks.Group("/get-files", Middleware.AuthAndRefreshMiddleware())
		{
			getFiles.GET("/", h.getFilesPage)
			getFiles.POST("/", h.getFiles)
		}
	}
}

func (h Handler) handleSubjectsPage(c *gin.Context, subj map[string]subjects) {
	if err := h.tmpl.ExecuteTemplate(c.Writer, "subjects.html", subj); err != nil {
		c.AbortWithStatus(InternalServerError)
		return
	}
}

func (h Handler) handleModulePage(c *gin.Context) {
	if err := h.tmpl.ExecuteTemplate(c.Writer, "module.html", nil); err != nil {
		c.AbortWithStatus(InternalServerError)
		return
	}
}

// StartServer runs the server
func (h Handler) StartServer() {
	h.handleUsers()
	h.handleLeaks()
	server := &http.Server{
		Addr:    ":8080",
		Handler: h.router,
		TLSConfig: &tls.Config{
			// Load the SSL/TLS certificate and private key files
			Certificates: []tls.Certificate{},
		},
	}
	err := server.ListenAndServeTLS("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
