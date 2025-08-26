package dependencies

import (
	userUsecase "cpi-hub-api/internal/core/usecase/user"
	userRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/json/user"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/user"
)

type Handlers struct {
	UserHandler *user.Handler
}

func Build() *Handlers {
	// repositories
	userRepository := userRepository.NewUserRepository("user_db.json")
	// genero el reposotory de spaces

	// usecases
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	// genero el usecase de spaces

	// handlers
	return &Handlers{
		UserHandler: &user.Handler{
			UseCase: userUsecase,
		},
		// SpaceHandler:....
	}
}
