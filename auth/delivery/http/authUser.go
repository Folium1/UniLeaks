package delivery

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (h Handler) AuthorizeUser(c *gin.Context, userId int) {
	refrToken := h.createRefreshToken(userId)
	authToken := h.createAuthToken(userId)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := h.useCase.SaveToken(ctx, refrToken)
	if err != nil {
		log.Println(err)
	}
	err = h.useCase.SaveToken(ctx, authToken)
	if err != nil {
		log.Println(err)
	}
	h.SetTokenToCookies(c, refrToken)
	h.SetTokenToCookies(c, authToken)
}

func (h Handler) LogOut(c *gin.Context) {
	h.deleteCookies(c)
}
