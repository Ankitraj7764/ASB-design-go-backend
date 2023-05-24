package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/design/responses"
)

type Credentials struct {
	Username string `json:"emailid"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}

/*
Login takes two things as parameter
- username or email
- password
it checks if the username and password are correct or not
if not it return invalid error

	else

it creates new jwt-access-token and refresh-token
and save it to the response body
*/
func Login(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("inside login function")
	//request body will constain user name and password
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("login req body decoded")

	// fmt.Println(credentials.Username)
	actualPassword := FindPassword(credentials.Username)
	fmt.Println(actualPassword)

	plainPwd := []byte(credentials.Password)
	if !ComparePasswords(actualPassword, plainPwd) {
		fmt.Println("password mismatched")
		rw.WriteHeader(http.StatusNonAuthoritativeInfo)
		return
	}
	fmt.Printf("password matched")

	accessToken, encryptedRefreshToken, AuthErr := AssignAuthTokens(credentials.Username)
	var loginResponse LoginResponse
	loginResponse.AccessToken = accessToken
	loginResponse.RefreshToken = encryptedRefreshToken
	if AuthErr != nil {
		rw.WriteHeader(http.StatusExpectationFailed)
		response := responses.UserResponse{Status: http.StatusExpectationFailed, Message: "unable to assign tokens", Data: map[string]interface{}{"data": AuthErr.Error()}}
		json.NewEncoder(rw).Encode(response)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": loginResponse}}
	json.NewEncoder(rw).Encode(response)

}
