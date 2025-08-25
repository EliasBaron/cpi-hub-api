package dependencies

import (
	userUsecase "cpi-hub-api/internal/core/usecase/user"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/user"
)

type Handlers struct {
	UserHandler *user.Handler
}

func Build() *Handlers {
	return &Handlers{
		UserHandler: &user.Handler{
			UseCase: userUsecase.NewUserUsecase(),
		},
	}
}
