package auth

import (
	"encoding/json"
	"net/http"

	"example.com/design/responses"
	"example.com/design/validators"
)

// var jwtKey = []byte("SecretKey")

func Home(rw http.ResponseWriter, r *http.Request) {
	// cookie, err := r.Cookie("token")
	// if err != nil {
	// 	if err == http.ErrNoCookie {
	// 		rw.WriteHeader(http.StatusUnauthorized)
	// 		return
	// 	}
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// tokenStr := cookie.Value
	// claims := &models.Claims{}
	// jwt_key_home, err := jwt.ParseWithClaims(tokenStr, claims,
	// 	func(t *jwt.Token) (interface{}, error) {
	// 		return jwtKey, nil
	// 	})
	// if err != nil {
	// 	if err == jwt.ErrSignatureInvalid {
	// 		rw.WriteHeader(http.StatusUnauthorized)
	// 		return
	// 	}
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// if !jwt_key_home.Valid {
	// 	rw.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// rw.Write([]byte(fmt.Sprintf("Authorized user :  %s", claims.Username)))

	//here I want validate the access token
	//return error if not validated
	//extract the token from the header
	//pass the token to the function validateAccessToken
	bearer := r.Header.Get("validator")
	err := validators.ValidateAccessToken(bearer)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(rw).Encode(response)
		return
	}

}
