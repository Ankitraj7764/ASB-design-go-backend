package auth

import (
	"time"

	"example.com/design/controllers"
)

func CreateNewReferenceToken(length int) string {
	//create random code
	randomCode := controllers.EncodeToString(length)
	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")
	randomUniqueCode := randomCode + currentTimeString
	// fmt.Println(randomUniqueCode)
	return randomUniqueCode
}
