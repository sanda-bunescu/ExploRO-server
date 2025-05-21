package repositories

import (
	"github.com/sanda-bunescu/ExploRO/models"
	"gorm.io/gorm"
)

type TouristicAttractionRepositoryInterface interface {
	GetAttractionsByCityId(cityId uint) ([]*models.TouristicAttraction, error)
	GetAttractionsThatAreNotInItineraryStopPoint(cityId uint, tripPlanId uint) ([]*models.TouristicAttraction, error)
	BaseRepositoryInterface[models.TouristicAttraction]
}

type TouristicAttractionRepository struct {
	BaseRepository[models.TouristicAttraction]
}

func NewTouristicAttractionRepository(db *gorm.DB) *TouristicAttractionRepository {
	return &TouristicAttractionRepository{BaseRepository: BaseRepository[models.TouristicAttraction]{DB: db}}
}

var _ TouristicAttractionRepositoryInterface = (*TouristicAttractionRepository)(nil)

func (t *TouristicAttractionRepository) GetAttractionsByCityId(cityId uint) ([]*models.TouristicAttraction, error) {
	var touristicAttractions []*models.TouristicAttraction
	err := t.DB.Where("city_id = ?", cityId).Find(&touristicAttractions).Error
	if err != nil {
		return nil, err
	}
	return touristicAttractions, nil
}

func (t *TouristicAttractionRepository) GetAttractionsThatAreNotInItineraryStopPoint(cityId uint, tripPlanId uint) ([]*models.TouristicAttraction, error) {
	var touristicAttractions []*models.TouristicAttraction
	err := t.DB.
		Table("touristic_attractions").
		Where("city_id = ? AND id NOT IN (?)", cityId,
			t.DB.Table("stop_points").
				Select("attraction_id").
				Joins("JOIN itineraries i ON stop_points.itinerary_id = i.id").
				Joins("JOIN trip_plans tp ON i.trip_plan_id = tp.id").
				Where("tp.id = ? AND stop_points.deleted_at IS NULL", tripPlanId),
		).
		Find(&touristicAttractions).Error

	if err != nil {
		return nil, err
	}
	return touristicAttractions, nil
}
