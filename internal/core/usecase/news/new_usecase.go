package news

import (
	"context"
	"cpi-hub-api/internal/core/domain"
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
	newsItems, err := u.newsRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return newsItems, nil
}

func (u *newsUsecase) CreateNews(ctx context.Context, new domain.News) (*domain.News, error) {
	createdNews, err := u.newsRepository.Create(ctx, &new)
	if err != nil {
		return nil, err
	}

	return createdNews, nil
}
