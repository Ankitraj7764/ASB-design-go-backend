package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id                      primitive.ObjectID `json:"id,omitempty"`
	Name                    string             `json:"name,omitempty" validate:"required"`
	EmailId                 string             `json:"email-id,omitempty" validate:"required" `
	Password                string             `json:"password,omitempty" validate:"required"`
	ProfilePictureURL       string             `json:"profile-picture-url,omitempty" validate:"required"`
	Score                   int                `json:"score,omitempty" validate:"required"`
	RecoveryCode            string             `json:"recovery-code"`
	CodeSendingTime         time.Time          `json:"code-sending-time"`
	SolvedChallenges        []Challange        `json:"solved-challenges"`
	RefreshTokens           []string           `json:"refresh-tokens"`
	CurrentAccessToken      string             `json:"current-access-token"`
	AccessTokenSendingTime  time.Time          `json:"access-token-sending-time"`
	RefreshTokenSendingTime time.Time          `json:"refresh-token-sending-time"`
	// SubmittedSolutions  []string
	// SubmittedChallenges []Challange
}

// coll := db.Collection("users")
