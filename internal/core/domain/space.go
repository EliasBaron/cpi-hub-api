package domain

import "time"

type Space struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	CreatedBy   int
	UpdatedAt   time.Time
	UpdatedBy   int
}

type SpaceWithUserAndCounts struct {
	Space       *Space
	User        *User
	SpaceCounts SpaceCounts
}

type SpaceCounts struct {
	Users int
	Posts int
}

type SpaceSearchCriteria struct {
	Name          *string
	CreatedBy     *int
	OrderBy       string
	Page          int
	PageSize      int
	SortDirection string
}

type SearchResult struct {
	Data     []*SpaceWithUserAndCounts
	Page     int
	PageSize int
	Total    int
}
