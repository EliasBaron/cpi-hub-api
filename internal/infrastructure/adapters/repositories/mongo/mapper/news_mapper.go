package mapper

import (
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToMongoNews(news *domain.News) *entity.News {
	var oid primitive.ObjectID
	if news.ID != "" {
		if parsed, err := primitive.ObjectIDFromHex(news.ID); err == nil {
			oid = parsed
		}
	}

	var expiresAt *primitive.DateTime
	if news.ExpiresAt != nil {
		dt := primitive.NewDateTimeFromTime(*news.ExpiresAt)
		expiresAt = &dt
	}

	return &entity.News{
		ID:          oid,
		Content:     news.Content,
		Image:       news.Image,
		RedirectURL: news.RedirectURL,
		ExpiresAt:   expiresAt,
		CreatedAt:   primitive.NewDateTimeFromTime(news.CreatedAt),
	}
}

func ToDomainNews(newsEntity *entity.News) *domain.News {
	var idStr string
	if newsEntity.ID != primitive.NilObjectID {
		idStr = newsEntity.ID.Hex()
	}
	var expiresAt *time.Time
	if newsEntity.ExpiresAt != nil {
		t := newsEntity.ExpiresAt.Time()
		expiresAt = &t
	}

	return &domain.News{
		ID:          idStr,
		Content:     newsEntity.Content,
		Image:       newsEntity.Image,
		RedirectURL: newsEntity.RedirectURL,
		ExpiresAt:   expiresAt,
		CreatedAt:   newsEntity.CreatedAt.Time(),
	}
}
