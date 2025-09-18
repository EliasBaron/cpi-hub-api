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
	Members     int
	Posts       int
}

type SpaceWithUser struct {
	Space *Space
	User  *User
}

type SpaceSearchCriteria struct {
	CreatedBy     *int
	OrderBy       string
	Page          int
	PageSize      int
	SortDirection string
}

type SearchResult struct {
	Data     []*SpaceWithUser
	Page     int
	PageSize int
	Total    int
}
