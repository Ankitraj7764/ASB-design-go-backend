package auth

import "example.com/design/controllers"

func FindPassword(userEmailId string) string {
	demoUser := controllers.FindUserByEmail(userEmailId)
	return demoUser.Password

}
