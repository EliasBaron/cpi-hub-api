package dto

type CreateSpace struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	CreatedBy   string `json:"created_by" binding:"required"`
}
