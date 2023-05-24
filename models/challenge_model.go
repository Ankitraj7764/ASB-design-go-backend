package models

// import "go.mongodb.org/mongo-driver/bson/primitive"
type Challange struct {
	ChallengeName        string `json:"challenge-name,omitempty" validate:"required"`
	Description          string `json:"description,omitempty" validate:"required"`
	ImageUrl             string `json:"image-url,omitempty" validate:"required"`
	Score                int    `json:"score-assigned,omitempty" validate:"required"`
	NextChallengeUrl     string `json:"next-challenge-url,omitempty" validate:"required"`
	PreviousChallengeUrl string `json:"previous-challenge-url,omitempty" validate:"required"`
	Difficulty           string `json:"difficulty,omitempty" validate:"required"`
	// easy/medium/hard
}
