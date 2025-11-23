package news

import (
	"context"
	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsRepository struct {
	db *mongo.Database
}

func NewNewsRepository(db *mongo.Database) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) GetAll(ctx context.Context, crit *criteria.Criteria) ([]*domain.News, error) {

	collection := r.db.Collection("news")

	filter := bson.M{}
	if crit != nil && len(crit.Filters) > 0 {
		criteriaFilter := mapper.ToMongoDBQuery(crit)
		for _, elem := range criteriaFilter {
			filter[elem.Key] = elem.Value
		}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	newsItems := make([]*domain.News, 0)
	for cursor.Next(ctx) {
		var newsEntity entity.News
		if err := cursor.Decode(&newsEntity); err != nil {
			return nil, err
		}
		news := mapper.ToDomainNews(&newsEntity)
		newsItems = append(newsItems, news)
	}
	return newsItems, nil
}

func (r *NewsRepository) Create(ctx context.Context, news *domain.News) (*domain.News, error) {
	newsEntity := mapper.ToMongoNews(news)

	collection := r.db.Collection("news")

	res, err := collection.InsertOne(ctx, newsEntity)
	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		news.ID = oid.Hex()
	}

	return news, nil
}
