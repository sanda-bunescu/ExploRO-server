package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type UserCityServiceInterface interface {
	AddUserCity(ginCtx *gin.Context, cityID uint) error
	DeleteUserCity(ginCtx *gin.Context, cityID uint) error
	DeleteUserCitiesByUserId(ctx context.Context, userId string) error
}

type UserCityService struct {
	UserCityRepo    repositories.UserCityRepositoryInterface
	CityRepo        repositories.CityRepositoryInterface
	firebaseService FirebaseServiceInterface
}

func NewUserCityService(UserCityRepo repositories.UserCityRepositoryInterface, CityRepo repositories.CityRepositoryInterface, firebaseService FirebaseServiceInterface) *UserCityService {
	return &UserCityService{UserCityRepo: UserCityRepo,
		CityRepo:        CityRepo,
		firebaseService: firebaseService}
}

var _ UserCityServiceInterface = (*UserCityService)(nil)

func (u *UserCityService) AddUserCity(ginCtx *gin.Context, cityID uint) error {
	//get UID
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return fmt.Errorf("unauthorized user: no firebaseUID in context")
	}
	//Check if city exists

	city, err := u.CityRepo.GetCityById(cityID) //getByCityId
	if err == nil && city == nil {
		return fmt.Errorf("city not found: %v", err)
	}
	if err != nil {
		return fmt.Errorf("error getting userCity: %v", err)
	}
	//Verify if user already has the city
	existingUserCity, err := u.UserCityRepo.GetUserCity(firebaseUID.(string), cityID)
	if err == nil && existingUserCity != nil {
		return fmt.Errorf("user already has this city added")
	}
	//Add city by user
	userCity := models.UserCities{
		UserID: firebaseUID.(string),
		CityID: cityID,
	}
	err = u.UserCityRepo.Create(ginCtx, &userCity)
	if err != nil {
		return fmt.Errorf("failed to add city: %v", err)
	}
	return nil
}

func (u *UserCityService) DeleteUserCity(ginCtx *gin.Context, cityID uint) error {
	//get UID
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return fmt.Errorf("unauthorized user: no firebaseUID in context")
	}

	foundUserCity, err := u.UserCityRepo.GetUserCity(firebaseUID.(string), cityID)
	if err != nil || foundUserCity == nil {
		return fmt.Errorf("GetUserCity: %w", err)
	}

	err = u.UserCityRepo.SoftDelete(ginCtx, foundUserCity)
	if err != nil {
		return fmt.Errorf("an error occured while executing delete %w", err)
	}
	return nil
}

func (u *UserCityService) DeleteUserCitiesByUserId(ctx context.Context, userId string) error {
	// Fetch all user cities
	userCities, err := u.UserCityRepo.GetUserCitiesByUserID(userId)
	if err != nil {
		return fmt.Errorf("failed to retrieve user cities: %w", err)
	}

	//No userCities found
	if len(userCities) == 0 {
		return nil
	}

	// Delete all user cities
	for _, city := range userCities {
		err := u.UserCityRepo.SoftDelete(ctx, city)
		if err != nil {
			return fmt.Errorf("failed to delete city %d for user: %w", city.CityID, err)
		}
	}

	return nil
}
