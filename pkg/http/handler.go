package delivery

import (
	"crypto/tls"
	"fmt"
	"html/template"
	auth "leaks/pkg/http/auth"
	admin "leaks/pkg/http/admin"
	leaksHandler "leaks/pkg/http/leaks"
	user "leaks/pkg/http/user"
	"leaks/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	logg       = logger.NewLogger()
	middleware = auth.New()
)

type Handler struct {
	tmpl   *template.Template
	leaks  *leaksHandler.LeaksHandler
	user   *user.UserHandler
	admin  *admin.AdminHandler
	router *gin.Engine
}

// New returns a new instance of handler
func New() Handler {
	tmpl := template.Must(template.ParseGlob("templates/*"))
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	// Load and parse the header template
	headerTemplate := template.Must(template.ParseFiles("templates/header.html"))
	// Set the template engine to use the parsed templates
	r.SetHTMLTemplate(headerTemplate)
	newUserHandler := user.New(tmpl)
	leaksHandler := leaksHandler.New(tmpl)
	adminHandler := admin.New(tmpl)
	return Handler{tmpl: tmpl, leaks: leaksHandler, user: newUserHandler, admin: adminHandler, router: r}
}

// handleUsers handles the user routes
func (h *Handler) handleUsers() {
	userGroup := h.router.Group("/user")
	{
		userGroup.GET("/register", h.user.Register)
		userGroup.GET("/login", h.user.Login)
		userGroup.POST("/register", h.user.PostRegister)
		userGroup.POST("/login", h.user.PostLogin)
		userGroup.GET("/logOut", h.user.LogOut)
	}
}

// handleLeaks handles the leaks routes
func (h *Handler) handleLeaks() {
	leaks := h.router.Group("/leaks", middleware.AuthAndRefreshMiddleware())
	{
		leaks.GET("/", h.leaks.MainPage)
		leaks.GET("/terms", h.termsOfUse)
		leaks.POST("/like-dislike", h.leaks.LikeDislikeFile)
		download := leaks.Group("/download")
		{
			download.GET("/:id", h.leaks.DownloadFile)
		}
		uploadFiles := leaks.Group("/")
		{
			uploadFiles.GET("/upload-files", h.leaks.UploadFilePage)
			uploadFiles.POST("/upload-files", h.leaks.UploadFile)
		}
		getFiles := leaks.Group("/get-files")
		{
			getFiles.GET("/", h.leaks.FilesPage)
			getFiles.POST("/", h.leaks.FilesList)
			getFiles.GET("/all", h.leaks.AllFiles)
		}
		myFiles := leaks.Group("/my-files")
		{
			myFiles.GET("/", h.leaks.MyFiles)
		}
	}
}

// handleAdmin handles the admin routes
func (h *Handler) handleAdmin() {
	admin := h.router.Group("/admin", middleware.AuthAndRefreshMiddleware(), middleware.OnlyAdminMiddleware())
	{
		admin.GET("/", h.admin.MainPage)
		leaks := admin.Group("/leaks")
		{
			leaks.GET("/files", h.admin.FilesList)
			leaks.GET("/file/:id", h.admin.DownloadFile)
			leaks.DELETE("/file/:id", h.admin.DeleteFile)

		}
		users := admin.Group("/users")
		{
			users.GET("/all", h.admin.AllUsers)
			users.POST("/ban", h.admin.BanUser)
			users.GET("/banned-users", h.admin.GetBannedUsers)
			users.POST("/unban/", h.admin.UnbanUser)
		}
	}
}

// termsOfUse handles the terms of use page
func (h *Handler) termsOfUse(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, "TermsOfUse.html", nil)
	if err != nil {
		logg.Error(err.Error())
	}
}

// StartServer runs the server
func (h *Handler) StartServer() {
	h.handleUsers()
	h.handleLeaks()
	h.handleAdmin()

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
		logg.Fatal(fmt.Sprint("ListenAndServeTLS: ", err))
	}
}
