package dependencies

import (
	postUsecase "cpi-hub-api/internal/core/usecase/post"
	spaceUsecase "cpi-hub-api/internal/core/usecase/space"
	userUsecase "cpi-hub-api/internal/core/usecase/user"
	postRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/post"
	spaceRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/space"
	userRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/user"
	userSpaceRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/user_space"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/post"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/space"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/user"
	"log"
)

type Handlers struct {
	UserHandler  *user.Handler
	SpaceHandler *space.SpaceHandler
	PostHandler  *post.PostHandler
}

func Build() *Handlers {

	sqldb, err := GetPostgreSQLDatabase()
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}

	userRepository := userRepository.NewUserRepository(sqldb)
	spaceRepository := spaceRepository.NewSpaceRepository(sqldb)
	userSpaceRepository := userSpaceRepository.NewUserSpaceRepository(sqldb)
	postRepository := postRepository.NewPostRepository(sqldb)

	userUsecase := userUsecase.NewUserUsecase(userRepository, spaceRepository, userSpaceRepository)
	spaceUsecase := spaceUsecase.NewSpaceUsecase(spaceRepository, userRepository, userSpaceRepository)
	postUsecase := postUsecase.NewPostUsecase(postRepository, spaceRepository, userRepository)

	return &Handlers{
		UserHandler: &user.Handler{
			UseCase: userUsecase,
		},
		SpaceHandler: &space.SpaceHandler{
			SpaceUseCase: spaceUsecase,
		},
		PostHandler: &post.PostHandler{
			PostUseCase: postUsecase,
		},
	}
}
