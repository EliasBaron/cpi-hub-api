package dto

import (
	"cpi-hub-api/internal/core/domain"
	"time"
)

type CreateNewsDTO struct {
	Content     string     `json:"content" binding:"required"`
	Image       string     `json:"image,omitempty"`
	RedirectURL string     `json:"redirect_url,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

func (dto *CreateNewsDTO) ToDomain() domain.News {
	news := domain.News{
		Content:     dto.Content,
		Image:       dto.Image,
		ExpiresAt:   dto.ExpiresAt,
		RedirectURL: dto.RedirectURL,
	}
	return news
}
