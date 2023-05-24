package validators

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"time"

	"example.com/design/configs"
	"example.com/design/controllers"
	"example.com/design/models"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// E represents a BSON element for a D. It is usually used inside a D
type E struct {
	Key   string
	Value interface{}
}

var currentUser models.User
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

/*
validate access token takes jwt access token as argument
and checks if access token is valid or noot
if not valid then it returns an error
else
it returns nil
*/
func ValidateAccessToken(requestAccessToken string) error {

	//validate access token
	//find user by email
	// currentUser = FindUserByEmail(email)
	// //validate access token
	// userAccessToken := currentUser.CurrentAccessToken
	// requestAccessTokenByte := []byte(requestAccessToken)
	// if !CompareTokens(userAccessToken, requestAccessTokenByte) {
	// 	return errors.New("invalid access token")
	// }
	//extract access token updated time
	// accessTokenSendingTime := currentUser.AccessTokenSendingTime
	// get current time
	// currentTime := time.Now()
	// elapsed := currentTime.Sub(accessTokenSendingTime)
	// //check if access token is expired or not
	// if elapsed <= 15*time.Minute {
	// 	//if not expired return nil
	// 	return nil
	// }
	//if expired get a new token
	// err := GetNewAccessToken(email)
	claims := &models.Claims{}
	secretKey := []byte("random@secret%key")
	jwt_key, err := jwt.ParseWithClaims(requestAccessToken, claims,
		func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

	if err != nil {
		return err
	}
	if !jwt_key.Valid {

		return errors.New("invalid access token")
	}

	return nil
}

/*
GetNewAccessTokenByRefreshTken takes two input as argument
- access token
- refresh token

if refresh token is expired it ask the user to log in again
else
it createes new refresh token and access token
and send it as response
*/
func GetNewAccessTokenByRefreshToken(oldRefreshToken, email string) (string, string, error) {

	// extract email from payload of old access token
	// var email string
	//get new user using email

	currentUser = FindUserByEmail(email)
	allRefreshTokens := currentUser.RefreshTokens
	TokenArraySize := len(allRefreshTokens)
	if TokenArraySize == 0 {
		return "", "", errors.New("no refresh token in the array")
	}

	//find the actual refresh token
	actualRefreshToken := currentUser.RefreshTokens[TokenArraySize-1]

	if !CompareTokens(actualRefreshToken, []byte(oldRefreshToken)) {

		for index := 0; index < TokenArraySize; index++ {
			//if we have any previous record for the refresh token that means invalidate all the refresh tokens
			if CompareTokens(actualRefreshToken, []byte(oldRefreshToken)) {
				//terminate all the previous tokens

				UpdateUserBody(email, nil)
				return "", "", errors.New("account is at risk! log in again")
				//return error and make the user login again
			}
		}
		//if all the refresh token do not match with request-refresh-token we are not going to do any
		return "", "", errors.New("invalid refresh token try again with right refresh token")
	}

	//get refresh token sent time

	refreshTokenSentTime := currentUser.RefreshTokenSendingTime

	//check if refresh token is expired or not

	//get current time

	currentTime := time.Now()

	elapsed := currentTime.Sub(refreshTokenSentTime)

	//check if access token is expired or not

	if elapsed <= 10*time.Hour {

		//if refresh token is valid generate a new access token and update access-token-sending time with current time

		newAccessToken, newRefreshToken, err := AssignAuthTokens(email)
		if err != nil {
			return "null", "null", err
		}
		return newAccessToken, newRefreshToken, nil

	}

	//if refresh token expired

	return "", "", errors.New("refresh token expired log in again")
}

/**********************************************************************************/
func AssignAuthTokens(email string) (string, string, error) {

	//generate access token

	newAccessToken, err := CreateNewJWTAccessToken(email)
	if err != nil {
		return "", "", err
	}
	// newAccessTokenByte := []byte(newAccessToken)
	// fmt.Println(newAccessToken)
	//generate refresh token

	newRefreshToken := CreateNewReferenceToken(6)
	newRefreshTokenByte := []byte(newRefreshToken)

	//encrypt access token

	// newEncryptedAccessToken := controllers.HashAndSalt(newAccessTokenByte)

	//encrypt refresh token

	newEncryptedRefreshToken := controllers.HashAndSalt(newRefreshTokenByte)

	//clear the refresh tokens array

	var newRefreshTokenArray []string

	//add new encrypted refresh token to the array

	newRefreshTokenArray = append(newRefreshTokenArray, newEncryptedRefreshToken)

	filter := bson.D{primitive.E{Key: "emailid", Value: email}}
	update := bson.D{{"$set", bson.D{primitive.E{Key: "refreshtokens", Value: newRefreshTokenArray}, primitive.E{Key: "refreshtokensendingtime", Value: time.Now()}}}}

	_, err2 := userCollection.UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func GenerateNewAccessToken(email string) error {

	//create random code

	// encryptedAccessToken := GenerateEncryptedToken(8)

	//encrypt new refresh token

	encryptedRefreshToken := GenerateEncryptedToken(6)

	//append refresh token to refresh-token-array

	oldRefreshTokenArray := currentUser.RefreshTokens

	//update refresh-token-array

	newRefreshTokenArray := append(oldRefreshTokenArray, encryptedRefreshToken)

	err := UpdateUserBody(email, newRefreshTokenArray)

	return err
}

func UpdateUserBody(email string, newRefreshTokenArray []string) error {

	//add new encrypted refresh token to the array

	filter := bson.D{{"emailid", email}}

	update := bson.D{{"$set", bson.D{
		//update refresh token array

		{"refreshtokens", newRefreshTokenArray},

		//update refresh token sent time

		{"refreshtokensendingtime", time.Now()},

		//update access token sent time
		// {"accesstokensendingtime", time.Now()},
		//update access token for that user
		// {"currentaccesstoken", encryptedAccessToken}
	}}}

	_, err := userCollection.UpdateOne(context.TODO(), filter, update)

	return err
}
func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}
func FindUserByEmail(targetEmailId string) models.User {

	coll := configs.DB.Database("golangAPI").Collection("users")
	filter := bson.D{{"emailid", targetEmailId}}

	var result models.User
	coll.FindOne(context.TODO(), filter).Decode(&result)

	return result
}
func GenerateEncryptedToken(length int) string {

	randomCode := EncodeToString(length)
	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")
	randomUniqueCode := randomCode + currentTimeString

	encryptedAccessToken := HashAndSalt([]byte(randomUniqueCode))
	return encryptedAccessToken

}

func CompareTokens(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func CreateNewReferenceToken(length int) string {
	//create random code
	randomCode := controllers.EncodeToString(length)
	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")
	randomUniqueCode := randomCode + currentTimeString
	// fmt.Println(randomUniqueCode)
	return randomUniqueCode
}

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
