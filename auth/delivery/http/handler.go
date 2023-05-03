package delivery

import "uniLeaks/auth"

type Handler struct {
	useCase auth.UseCase
}

func New() Handler {
	return Handler{auth.NewUseCase()}
}
