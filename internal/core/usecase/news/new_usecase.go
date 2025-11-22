package news

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/pkg/helpers"
)

type NewsUseCase interface {
	GetAllNews(ctx context.Context) ([]*domain.News, error)
	CreateNews(ctx context.Context, new domain.News) (*domain.News, error)
}

type newsUsecase struct {
	newsRepository domain.NewsRepository
}

func NewNewsUsecase(newsRepo domain.NewsRepository) NewsUseCase {
	return &newsUsecase{
		newsRepository: newsRepo,
	}
}

func (u *newsUsecase) GetAllNews(ctx context.Context) ([]*domain.News, error) {

	now := helpers.GetTime()
	searchCriteria := criteria.NewCriteriaBuilder().
		WithFilter("expires_at", now, criteria.OperatorGte).
		Build()

	newsItems, err := u.newsRepository.GetAll(ctx, searchCriteria)
	if err != nil {
		return nil, err
	}

	return newsItems, nil
}

func (u *newsUsecase) CreateNews(ctx context.Context, new domain.News) (*domain.News, error) {

	new.CreatedAt = helpers.GetTime()

	//si no se indicó fecha de expiración, se asigna 1 mes desde la creación
	if new.ExpiresAt == nil {
		t := new.CreatedAt.AddDate(0, 1, 0)
		new.ExpiresAt = &t
	}

	createdNews, err := u.newsRepository.Create(ctx, &new)
	if err != nil {
		return nil, err
	}

	return createdNews, nil
}
