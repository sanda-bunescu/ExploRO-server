package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/cmd/container"
	"github.com/sanda-bunescu/ExploRO/cmd/middleware"
	"github.com/sanda-bunescu/ExploRO/controllers"
)

func RegisterRoutes(router *gin.Engine, c *container.Container) {
	StaticRoutes(router)
	TouristicAttractionController(router, c)
	UserRoutes(router, c)
	UserCityRoutes(router, c)
	UserGroupsRoutes(router, c)
	GroupRoutes(router, c)
	TripPlanRoutes(router, c)
	ItineraryRoutes(router, c)
	StopPointRoutes(router, c)
	ExpenseRoutes(router, c)
	DebtRoutes(router, c)
}

func StaticRoutes(router *gin.Engine) {
	router.Static("/static", "./static")
}

func UserRoutes(router *gin.Engine, c *container.Container) {
	router.Use(middleware.FirebaseAuthMiddleware(c.FirebaseService))

	userController := controllers.NewUserController(c.UserService, c.FirebaseService, c.UserCityService)
	router.POST("/create-user", userController.Register)
	router.POST("/login-user", userController.Login)
	router.POST("/delete-user", userController.DeleteUser)
	router.GET("/get-user", userController.GetUserByID)
	router.GET("/get-user-cities", userController.GetUserCitiesByUserID)
	router.POST("/add-user-city", userController.AddUserCity)
	router.DELETE("/delete-user-city", userController.DeleteUserCity)
}

func TouristicAttractionController(router *gin.Engine, c *container.Container) {
	touristicAttractionController := controllers.NewTouristicAttractionController(c.TouristicAttractionService)
	router.GET("/get-attractions", touristicAttractionController.GetAllTouristicAttractions)
	router.GET("/get-attractions-by-cityid", touristicAttractionController.GetAllByCityId)
	router.GET("/attractions/not-in-itinerary", touristicAttractionController.GetAttractionsNotInItinerary)

}

func UserGroupsRoutes(router *gin.Engine, c *container.Container) {
	userGroupController := controllers.NewUserGroupController(c.UserGroupService)
	router.GET("/get-users-by-groupid", userGroupController.GetAllUsersByGroupId)
	router.GET("/get-groups-by-userid", userGroupController.GetAllByUserID)
	router.PUT("/add-user-groups", userGroupController.AddUserGroup)
	router.DELETE("/delete-user-groups", userGroupController.DeleteUserGroup)
}

func UserCityRoutes(router *gin.Engine, c *container.Container) {
	cityController := controllers.NewCityController(c.CityService)
	router.GET("/get-cities", cityController.GetAllCities)
}

func GroupRoutes(router *gin.Engine, c *container.Container) {
	groupController := controllers.NewGroupController(c.GroupService)
	router.POST("/create-group", groupController.CreateGroup)
	router.DELETE("/delete-group", groupController.DeleteGroup)
}

func TripPlanRoutes(router *gin.Engine, c *container.Container) {
	tripPlanController := controllers.NewTripPlanController(c.TripPlanService)
	router.POST("/create-trip", tripPlanController.CreateTrip)
	router.DELETE("/delete-trip", tripPlanController.DeleteTrip)
	router.GET("/get-trip-by-user", tripPlanController.GetTripsByUserId)
	router.GET("/get-trip-by-city-and-user", tripPlanController.GetTripsByCityAndUser)
	router.GET("/get-trips-by-group", tripPlanController.GetTripsByGroupId)
}

func ItineraryRoutes(router *gin.Engine, c *container.Container) {
	itineraryController := controllers.NewItineraryController(c.ItineraryService)
	router.GET("/get-all-itineraries-by-trip-plan-id", itineraryController.GetByTripPlanId)
	router.POST("/create-itinerary", itineraryController.Create)
	router.DELETE("/delete-itinerary", itineraryController.Delete)
}

func StopPointRoutes(router *gin.Engine, c *container.Container) {
	stopPointController := controllers.NewStopPointController(c.StopPointService)
	router.GET("/get-all-attractions-by-itinerary-id", stopPointController.GetAllByItineraryId)
	router.POST("/add-touristic-attractions-itinerary", stopPointController.Create)
	router.DELETE("/delete-touristic-attraction", stopPointController.Delete)
}

func ExpenseRoutes(router *gin.Engine, c *container.Container) {
	expenseController := controllers.NewExpenseController(c.ExpenseService)
	router.GET("/get-expenses-by-groupid", expenseController.GetAllExpensesByGroupId)
	router.POST("/save-expense", expenseController.CreateExpenseWithDebts)
	router.PUT("/edit-expense", expenseController.EditExpenseWithDebts)
	router.DELETE("/delete-expense", expenseController.DeleteExpenseByID)
}

func DebtRoutes(router *gin.Engine, c *container.Container) {
	debtController := controllers.NewDebtController(c.DebtService)
	router.GET("get-user-debts", debtController.GetByGroupIdAndUser)
	router.DELETE("delete-debt", debtController.DeleteDebt)
}
