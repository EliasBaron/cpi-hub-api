package dependencies

import (
	spaceUsecase "cpi-hub-api/internal/core/usecase/space"
	userUsecase "cpi-hub-api/internal/core/usecase/user"
	spaceRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/space"
	userRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgre/user"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/space"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/user"
	"log"
)

type Handlers struct {
	UserHandler  *user.Handler
	SpaceHandler *space.SpaceHandler
}

func Build() *Handlers {
	mongoDB, err := GetMongoDatabase()
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	postgreDB, err := GetPostgreSQLDatabase()
	if err != nil {
		log.Fatalf("Error al conectar a PostgreSQL: %v", err)
	}

	userRepository := userRepository.NewUserRepository(postgreDB)
	spaceRepository := spaceRepository.NewSpaceRepository(mongoDB)

	userUsecase := userUsecase.NewUserUsecase(userRepository)
	spaceUsecase := spaceUsecase.NewSpaceUsecase(spaceRepository, userRepository)

	return &Handlers{
		UserHandler: &user.Handler{
			UseCase: userUsecase,
		},
		SpaceHandler: &space.SpaceHandler{
			SpaceUseCase: spaceUsecase,
		},
	}
}
