package api

import (
	"crypto/tls"
	"fmt"
	"net/http"

	adminApi "leaks/pkg/admin/http"
	auth "leaks/pkg/auth"
	leaksApi "leaks/pkg/leaks/http"
	"leaks/pkg/logger"
	userApi "leaks/pkg/user/http"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	leaks  *leaksApi.LeaksHandler
	user   *userApi.UserHandler
	admin  *adminApi.AdminHandler
	router *gin.Engine
}

var (
	l          = logger.NewLogger()
	middleware = auth.New()
)

func NewApiHandler() ApiHandler {
	r := gin.Default()

	user := userApi.New()
	leaks := leaksApi.New()
	admin := adminApi.New()
	return ApiHandler{
		leaks:  leaks,
		user:   user,
		admin:  admin,
		router: r,
	}
}

func (h *ApiHandler) handleUsers() {
	userGroup := h.router.Group("/user")
	{
		userGroup.POST("/register", h.user.PostRegister)
		userGroup.POST("/login", h.user.PostLogin)
		userGroup.GET("/logOut", h.user.LogOut)
	}
}

func (h *ApiHandler) handleLeaks() {
	leaks := h.router.Group("/leaks", middleware.AuthAndRefreshMiddleware())
	{
		leaks.GET("/", h.leaks.MainPage)
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

func (h *ApiHandler) handleAdmin() {
	admin := h.router.Group("/admin", middleware.AuthAndRefreshMiddleware(), middleware.OnlyAdminMiddleware())
	{
		admin.GET("/", h.admin.MainPage)
		leaks := admin.Group("/leaks")
		{
			leaks.GET("/files", h.admin.FilesOrderedByDislikes)
			leaks.GET("/file/:id", h.admin.DownloadFile)
			leaks.DELETE("/file/:id", h.admin.DeleteFile)

		}
		users := admin.Group("/users")
		{
			users.GET("/all", h.admin.AllUsers)
			users.POST("/ban-user", h.admin.BanUser)
			users.GET("/banned-users", h.admin.GetBannedUsers)
			users.POST("/unban/", h.admin.UnbanUser)
		}
	}
}

func (h *ApiHandler) StartServer() {
	h.handleUsers()
	h.handleLeaks()
	h.handleAdmin()

	server := &http.Server{
		Addr:    ":8080",
		Handler: h.router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{},
		},
	}
	err := server.ListenAndServeTLS("certs/server.crt", "certs/server.key")
	if err != nil {
		l.Fatal(fmt.Sprint("ListenAndServeTLS: ", err))
	}
}
