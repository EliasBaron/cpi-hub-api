package dependencies

import (
	userUsecase "cpi-hub-api/internal/core/usecase/user"
	userRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/json"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/user"
)

type Handlers struct {
	UserHandler *user.Handler
}

func Build() *Handlers {
	userRepository := userRepository.NewUserRepository("user_db.json")
	return &Handlers{
		UserHandler: &user.Handler{
			UseCase: userUsecase.NewUserUsecase(userRepository),
		},
	}
}
