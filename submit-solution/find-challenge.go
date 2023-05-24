package submitsolution

import (
	"context"
	"fmt"

	"example.com/design/models"
	"go.mongodb.org/mongo-driver/bson"
)

func FindChallengeByName(name string) (models.Challange, error) {
	var targetChallenge models.Challange
	fmt.Println(name)
	filter := bson.D{{"challengename", name}}
	err := challengeCollection.FindOne(context.TODO(), filter).Decode(&targetChallenge)
	if err != nil {
		return targetChallenge, err
	}
	return targetChallenge, nil

}
