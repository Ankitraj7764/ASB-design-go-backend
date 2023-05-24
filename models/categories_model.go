package models

type Categories struct {
	CategoryName string `json:"category-name,omitempty" validate:"required"`
	ProblemList  []Challange
}
