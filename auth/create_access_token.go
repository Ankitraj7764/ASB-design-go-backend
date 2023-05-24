package auth

import (
	"time"

	"example.com/design/models"
	"github.com/golang-jwt/jwt"
)

//create new JWT token and {refresh token}
//response will contain the JWT token
//this will be stored in the ls
//in the header they should put this jwt token
//we will validate that and if expired we will create a new one
//

func CreateNewJWTAccessToken(email string) (string, error) {
	tokenLifetime := 10 * time.Hour
	secretKey := []byte("random@secret%key")
	claims := models.JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(tokenLifetime)).Unix(),
		},
		CustmClaims: map[string]string{
			"email": email,
		},
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println("access token-string is : ", tokenString)
	token, err := tokenString.SignedString(secretKey)
	// fmt.Println("access token is : ", token)
	// fmt.Println("error at line 33 : ", err)
	return token, err
	//create random code
	// randomCode := controllers.EncodeToString(length)
	// currentTime := time.Now()
	// currentTimeString := currentTime.Format("2006-01-02 15:04:05")
	// randomUniqueCode := randomCode + currentTimeString
	// fmt.Println(randomUniqueCode)
	// return randomUniqueCode

}
