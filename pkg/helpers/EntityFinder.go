package helpers

import (
	"context"
	"cpi-hub-api/internal/core/domain/criteria"
	"cpi-hub-api/pkg/apperror"
)

type EntityFinder[T any] interface {
	Find(ctx context.Context, c *criteria.Criteria) (*T, error)
}

func FindEntity[T any](
	ctx context.Context,
	repo EntityFinder[T],
	field string,
	value interface{},
	notFoundMsg string,
) (*T, error) {
	entity, err := repo.Find(ctx, &criteria.Criteria{
		Filters: []criteria.Filter{
			{Field: field, Value: value, Operator: criteria.OperatorEqual},
		},
	})
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, apperror.NewNotFound(notFoundMsg, nil, "EntityFinder.go:FindEntity")
	}
	return entity, nil
}
