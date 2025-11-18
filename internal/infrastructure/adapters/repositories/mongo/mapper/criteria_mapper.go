package mapper

import (
	"cpi-hub-api/internal/core/domain/criteria"

	"go.mongodb.org/mongo-driver/bson"
)

func ToMongoDBQuery(c *criteria.Criteria) bson.D {
	filters := bson.D{}
	for _, filter := range c.Filters {
		switch filter.Operator {
		case criteria.OperatorEqual:
			filters = append(filters, bson.E{Key: filter.Field, Value: filter.Value})
		case criteria.OperatorNotEqual:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$ne": filter.Value}})
		case criteria.OperatorIn:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$in": filter.Value}})
		case criteria.OperatorNotIn:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$nin": filter.Value}})
		case criteria.OperatorGt:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$gt": filter.Value}})
		case criteria.OperatorGte:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$gte": filter.Value}})
		case criteria.OperatorLt:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$lt": filter.Value}})
		case criteria.OperatorLte:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$lte": filter.Value}})
		case criteria.OperatorRegex:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$regex": filter.Value}})
		case criteria.OperatorExists:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$exists": filter.Value}})
		case criteria.OperatorILike:
			filters = append(filters, bson.E{Key: filter.Field, Value: bson.M{"$regex": filter.Value, "$options": "i"}})
		}
	}
	return filters
}
