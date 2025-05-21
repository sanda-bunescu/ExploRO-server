package repositories

import (
	"errors"
	"github.com/sanda-bunescu/ExploRO/models"
	"gorm.io/gorm"
)

type ItineraryRepositoryInterface interface {
	BaseRepositoryInterface[models.Itinerary]
	GetAllByTripPlanId(tripPlanId uint) ([]*models.Itinerary, error)
	GetLastAddedItinerary(tripPlanId uint) (*models.Itinerary, error)
}

type ItineraryRepository struct {
	BaseRepository[models.Itinerary]
}

func NewItineraryRepository(db *gorm.DB) *ItineraryRepository {
	return &ItineraryRepository{BaseRepository[models.Itinerary]{DB: db}}
}

var _ ItineraryRepositoryInterface = (*ItineraryRepository)(nil)

func (ir *ItineraryRepository) GetAllByTripPlanId(tripPlanId uint) ([]*models.Itinerary, error) {
	var itineraries []*models.Itinerary
	err := ir.DB.Where("trip_plan_id = ? AND deleted_at IS NULL", tripPlanId).Find(&itineraries).Error
	if err != nil {
		return nil, err
	}
	return itineraries, nil
}

func (ir *ItineraryRepository) GetLastAddedItinerary(tripPlanId uint) (*models.Itinerary, error) {
	var itineraries *models.Itinerary
	err := ir.DB.Where("trip_plan_id = ? AND deleted_at IS NULL", tripPlanId).Order("day_nr DESC").First(&itineraries).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return itineraries, nil
}
