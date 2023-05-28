package delivery

import (
	"uniLeaks/auth"
)

type Handler struct {
	useCase auth.UseCase
}

// New returns a new instance of the auth handler.
func New() Handler {
	return Handler{auth.NewUseCase()}
}
