package dependencies

import (
	spaceUsecase "cpi-hub-api/internal/core/usecase/space"
	userUsecase "cpi-hub-api/internal/core/usecase/user"
	spaceRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/space"
	userRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/user"
	userSpaceRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/user_space"
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
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	sqldb, err := GetPostgreSQLDatabase()
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}

	userRepository := userRepository.NewUserRepository(sqldb)
	spaceRepository := spaceRepository.NewSpaceRepository(mongoDB)
	userSpaceRepository := userSpaceRepository.NewUserSpaceRepository(sqldb)

	userUsecase := userUsecase.NewUserUsecase(userRepository, spaceRepository, userSpaceRepository)
	spaceUsecase := spaceUsecase.NewSpaceUsecase(spaceRepository, userRepository, userSpaceRepository)

	return &Handlers{
		UserHandler: &user.Handler{
			UseCase: userUsecase,
		},
		SpaceHandler: &space.SpaceHandler{
			SpaceUseCase: spaceUsecase,
		},
	}
}
