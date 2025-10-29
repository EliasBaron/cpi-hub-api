package entity

type Reaction struct {
	ID         string `bson:"id"`
	UserID     int    `bson:"user_id"`
	EntityType string `bson:"entity_type"`
	EntityID   int    `bson:"entity_id"`
	Liked      bool   `bson:"liked"`
	Disliked   bool   `bson:"disliked"`
}
