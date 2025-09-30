package dependencies

import (
	commentUsecase "cpi-hub-api/internal/core/usecase/comment"
	postUsecase "cpi-hub-api/internal/core/usecase/post"
	spaceUsecase "cpi-hub-api/internal/core/usecase/space"
	userUsecase "cpi-hub-api/internal/core/usecase/user"
	commentRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/comment"
	postRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/post"
	spaceRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/space"
	userRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/user"
	userSpaceRepository "cpi-hub-api/internal/infrastructure/adapters/repositories/postgres/user_space"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/comment"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/post"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/space"
	"cpi-hub-api/internal/infrastructure/entrypoint/handlers/user"
	"log"
)

type Handlers struct {
	UserHandler    *user.UserHandler
	SpaceHandler   *space.SpaceHandler
	PostHandler    *post.PostHandler
	CommentHandler *comment.CommentHandler
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
	commentRepository := commentRepository.NewCommentRepository(sqldb)

	userUsecase := userUsecase.NewUserUsecase(userRepository, spaceRepository, userSpaceRepository)
	spaceUsecase := spaceUsecase.NewSpaceUsecase(spaceRepository, userRepository, userSpaceRepository, postRepository)
	postUsecase := postUsecase.NewPostUsecase(postRepository, spaceRepository, userRepository, commentRepository, userSpaceRepository)
	commentUsecase := commentUsecase.NewCommentUsecase(commentRepository)

	return &Handlers{
		UserHandler: &user.UserHandler{
			UseCase:     userUsecase,
			PostUseCase: postUsecase,
		},
		SpaceHandler: &space.SpaceHandler{
			SpaceUseCase: spaceUsecase,
		},
		PostHandler: &post.PostHandler{
			PostUseCase: postUsecase,
		},
		CommentHandler: &comment.CommentHandler{
			CommentUseCase: commentUsecase,
		},
	}
}
