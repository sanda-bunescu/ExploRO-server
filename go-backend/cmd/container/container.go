package container

import (
	"github.com/sanda-bunescu/ExploRO/initializers"
	"github.com/sanda-bunescu/ExploRO/repositories"
	"github.com/sanda-bunescu/ExploRO/services"
)

type Container struct {
	UserService                *services.UserService
	UserCityService            *services.UserCityService
	FirebaseService            *services.FirebaseService
	CityService                *services.CityService
	TouristicAttractionService *services.TouristicAttractionService
	UserGroupService           *services.UserGroupService
	GroupService               *services.GroupService
	TripPlanService            *services.TripPlanService
	ItineraryService           *services.ItineraryService
	StopPointService           *services.StopPointService
	ExpenseService             *services.ExpenseService
	DebtService                *services.DebtService
}

func NewContainer() (*Container, error) {
	//initialize repos
	userRepo := repositories.NewUserRepository(initializers.Database)
	userCityRepo := repositories.NewUserCityRepository(initializers.Database)
	cityRepo := repositories.NewCityRepository(initializers.Database)
	touristicAttractionRepo := repositories.NewTouristicAttractionRepository(initializers.Database)
	userGroupRepo := repositories.NewUserGroupRepository(initializers.Database)
	groupRepo := repositories.NewGroupRepository(initializers.Database)
	tripPlanRepo := repositories.NewTripPlanRepository(initializers.Database)
	itineraryRepo := repositories.NewItineraryRepository(initializers.Database)
	stopPointRepo := repositories.NewStopPointRepository(initializers.Database)
	expenseRepo := repositories.NewExpenseRepository(initializers.Database)
	debtRepo := repositories.NewDebtRepository(initializers.Database)

	// Initialize Firebase service
	firebaseService, err := services.NewFirebaseService(userRepo)
	if err != nil {
		return nil, err
	}

	// Initialize services
	tripPlanService := services.NewTripPlanService(tripPlanRepo)
	userCityService := services.NewUserCityService(userCityRepo, cityRepo, firebaseService)
	debtService := services.NewDebtService(debtRepo, userRepo, expenseRepo)
	expenseService := services.NewExpenseService(expenseRepo, userRepo, debtService)
	groupService := services.NewGroupService(groupRepo, firebaseService, tripPlanService, userGroupRepo, expenseService)
	userService := services.NewUserService(userRepo, userCityService, groupService, firebaseService, expenseService)
	userGroupService := services.NewUserGroupService(userGroupRepo, groupService, userRepo)
	cityService := services.NewCityService(cityRepo)
	touristicAttractionService := services.NewTouristicAttractionService(touristicAttractionRepo)
	stopPointService := services.NewStopPointService(stopPointRepo)
	itineraryService := services.NewItineraryService(itineraryRepo, stopPointService)

	// Return a container with dependencies
	return &Container{
		UserService:                userService,
		UserCityService:            userCityService,
		CityService:                cityService,
		TouristicAttractionService: touristicAttractionService,
		FirebaseService:            firebaseService,
		UserGroupService:           userGroupService,
		GroupService:               groupService,
		TripPlanService:            tripPlanService,
		ItineraryService:           itineraryService,
		StopPointService:           stopPointService,
		ExpenseService:             expenseService,
		DebtService:                debtService,
	}, nil
}
