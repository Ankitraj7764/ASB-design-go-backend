package models

type Level struct {
	LevelName    string `json:"level-name,omitempty" validate:"required"`
	CategoryList []Categories
}
