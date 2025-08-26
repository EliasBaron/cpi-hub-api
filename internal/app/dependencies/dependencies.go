package dependencies

import (
	spaceUsecase "cpi-hub-api/internal/core/usecase/space"
	userUsecase "cpi-hub-api/internal/core/usecase/user"
	spaceRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/json/space"
	userRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/json/user"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/space"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/user"
)

type Handlers struct {
	UserHandler  *user.Handler
	SpaceHandler *space.SpaceHandler
}

func Build() *Handlers {
	// repositories
	userRepository := userRepository.NewUserRepository("user_db.json")
	// genero el reposotory de spaces
	spaceRepository := spaceRepository.NewSpaceRepository("space_db.json")

	// usecases
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	// genero el usecase de spaces
	spaceUsecase := spaceUsecase.NewSpaceUsecase(spaceRepository)

	// handlers
	return &Handlers{
		UserHandler: &user.Handler{
			UseCase: userUsecase,
		},
		SpaceHandler: &space.SpaceHandler{
			SpaceUseCase: spaceUsecase,
		},
	}
}
