package repositories

import (
	"github.com/sanda-bunescu/ExploRO/models"
	"gorm.io/gorm"
)

type StopPointRepositoryInterface interface {
	BaseRepositoryInterface[models.StopPoint]
	GetAllByItineraryId(itineraryId uint) ([]*models.StopPoint, error)
}

type StopPointRepository struct {
	BaseRepository[models.StopPoint]
}

func NewStopPointRepository(db *gorm.DB) *StopPointRepository {
	return &StopPointRepository{BaseRepository[models.StopPoint]{DB: db}}
}

var _ StopPointRepositoryInterface = (*StopPointRepository)(nil)

func (spr *StopPointRepository) GetAllByItineraryId(itineraryId uint) ([]*models.StopPoint, error) {
	var stopPoints []*models.StopPoint
	err := spr.DB.
		Preload("TouristicAttraction").
		Where("itinerary_id = ? AND deleted_at IS NULL", itineraryId).
		Order("itinerary_stop_nr ASC").
		Find(&stopPoints).Error
	if err != nil {
		return nil, err
	}
	return stopPoints, nil
}
