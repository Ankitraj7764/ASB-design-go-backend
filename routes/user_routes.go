package routes

import (
	"example.com/design/auth"
	"example.com/design/categories"
	"example.com/design/controllers"
	"example.com/design/levels"
	"example.com/design/modules"
	submitsolution "example.com/design/submit-solution"
	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {

	//user validation requests
	router.HandleFunc("/login", auth.Login)
	router.HandleFunc("/log-out", auth.LogOut)
	router.HandleFunc("/home", auth.Home)
	router.HandleFunc("/forgot-password", controllers.ForgotPassword)
	router.HandleFunc("/reset-password", controllers.ResetPassword)

	//create requests

	router.HandleFunc("/user", controllers.CreateUser()).Methods("POST")
	router.HandleFunc("/module", modules.CreateModule()).Methods("POST")
	router.HandleFunc("/level", levels.CreateLevel()).Methods("POST")
	// router.HandleFunc("/level", levels.CreateLevel()).Methods("POST")
	router.HandleFunc("/add-challenge", submitsolution.CreateChallenge()).Methods("POST")
	router.HandleFunc("/category", categories.CreateCategory()).Methods("POST")

	//read requests

	router.HandleFunc("/user/{userId}", controllers.GetAUser()).Methods("GET")
	router.HandleFunc("/get-all-users", controllers.GetAllUsers()).Methods("GET")
	router.HandleFunc("/get-all-submissions", submitsolution.GetAllSubmissions()).Methods("GET")
	router.HandleFunc("/get-leaderboard-details", controllers.GetLeaderboardDetails()).Methods("GET")
	router.HandleFunc("/find-module", modules.FindModule()).Methods("POST")
	router.HandleFunc("/find-level", levels.FindLevel()).Methods("POST")
	router.HandleFunc("/get-all-levels", levels.GetAllLevels()).Methods("GET")
	router.HandleFunc("/get-all-challenges", submitsolution.GetAllChallenges()).Methods("GET")
	router.HandleFunc("/get-all-categories", categories.GetAllCategories()).Methods("GET")
	router.HandleFunc("/get-all-modules", modules.GetAllModules()).Methods("GET")
	router.HandleFunc("/find-category", categories.FindCategory()).Methods("POST")
	// router.HandleFunc("/module/{moduleName}", modules.GetAModule()).Methods("GET")

	// update requests

	router.HandleFunc("/add-level-in-module", modules.AddLevelInModule()).Methods("POST")
	router.HandleFunc("/add-category-in-level", levels.AddCategoryInLevel()).Methods("POST")
	router.HandleFunc("/add-challenge-in-category", categories.AddChallengeInCategory()).Methods("POST")

	// delete requests

	router.HandleFunc("/delete-level", levels.DeleteLevel()).Methods("POST")
	router.HandleFunc("/delete-module", modules.DeleteModule()).Methods("POST")
	router.HandleFunc("/delete-category", categories.DeleteCategory()).Methods("POST")
	router.HandleFunc("/accept-solution", submitsolution.AcceptSolution()).Methods("POST")

	//special requests

	router.HandleFunc("/submit-solution", submitsolution.SubmitSolution()).Methods("POST")

}
