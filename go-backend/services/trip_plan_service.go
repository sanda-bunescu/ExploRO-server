package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type TripPlanServiceInterface interface {
	GetTripsByUserId(ginCtx *gin.Context) ([]*responses.TripPlanResponse, error)
	GetTripsByCityAndUser(cityId uint, ginCtx *gin.Context) ([]*responses.TripPlanResponse, error)
	GetTripsByGroupId(groupId uint) ([]*responses.TripPlanResponse, error)
	CreateTrip(ctx context.Context, trip *requests.CreateTripPlanRequest) error
	DeleteTrip(ctx context.Context, tripId uint) error
}

type TripPlanService struct {
	tripPlanRepository repositories.TripPlanRepositoryInterface
}

func NewTripPlanService(tripPlanRepo repositories.TripPlanRepositoryInterface) *TripPlanService {
	return &TripPlanService{tripPlanRepository: tripPlanRepo}
}

var _ TripPlanServiceInterface = (*TripPlanService)(nil)

func (t TripPlanService) GetTripsByUserId(ginCtx *gin.Context) ([]*responses.TripPlanResponse, error) {
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user: no firebaseUID in context")
	}

	response, err := t.tripPlanRepository.GetTripsByUserId(firebaseUID.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to get user's trips: %w", err)
	}

	return response, nil

}

func (t TripPlanService) GetTripsByCityAndUser(cityId uint, ginCtx *gin.Context) ([]*responses.TripPlanResponse, error) {
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user: no firebaseUID in context")
	}

	response, err := t.tripPlanRepository.GetTripsByCityAndUser(cityId, firebaseUID.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to get trips for city %d and user: %w", cityId, err)
	}

	return response, nil
}

func (t TripPlanService) GetTripsByGroupId(groupId uint) ([]*responses.TripPlanResponse, error) {
	response, err := t.tripPlanRepository.GetTripsByGroupId(groupId)
	if err != nil {
		return nil, fmt.Errorf("failed to get trips for group %d: %w", groupId, err)
	}

	return response, nil
}

func (t TripPlanService) CreateTrip(ctx context.Context, trip *requests.CreateTripPlanRequest) error {
	tripPlan := models.TripPlan{
		Name:      trip.TripName,
		StartDate: trip.StartDate,
		EndDate:   trip.EndDate,
		CityId:    trip.CityId,
		GroupId:   trip.GroupId,
	}

	err := t.tripPlanRepository.Create(ctx, &tripPlan)
	if err != nil {
		return fmt.Errorf("failed to create trip plan: %w", err)
	}

	return nil
}

func (t TripPlanService) DeleteTrip(ctx context.Context, tripId uint) error {
	trip, err := t.tripPlanRepository.GetByID(ctx, tripId)
	if err != nil {
		return fmt.Errorf("failed to find trip plan with ID %d: %w", tripId, err)
	}

	err = t.tripPlanRepository.SoftDelete(ctx, trip)
	if err != nil {
		return fmt.Errorf("failed to soft delete trip plan with ID %d: %w", tripId, err)
	}

	return nil
}
