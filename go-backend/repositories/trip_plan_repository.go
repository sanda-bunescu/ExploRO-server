package repositories

import (
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"gorm.io/gorm"
)

type TripPlanRepositoryInterface interface {
	GetTripsByUserId(firebaseUID string) ([]*responses.TripPlanResponse, error)
	GetTripsByCityAndUser(cityId uint, firebaseUID string) ([]*responses.TripPlanResponse, error)
	GetTripsByGroupId(groupId uint) ([]*responses.TripPlanResponse, error)

	BaseRepositoryInterface[models.TripPlan]
}

type TripPlanRepository struct {
	BaseRepository[models.TripPlan]
}

func NewTripPlanRepository(db *gorm.DB) *TripPlanRepository {
	return &TripPlanRepository{BaseRepository[models.TripPlan]{db}}
}

var _ TripPlanRepositoryInterface = (*TripPlanRepository)(nil)

func (r *TripPlanRepository) GetTripsByUserId(firebaseUID string) ([]*responses.TripPlanResponse, error) {
	var tripPlans []*responses.TripPlanResponse

	err := r.DB.
		Table("trip_plans").
		Select("trip_plans.id AS id, trip_plans.name AS trip_name, trip_plans.start_date AS start_date, trip_plans.end_date AS end_date, groups.name AS group_name, cities.name AS city_name, cities.id as city_id").
		Joins("JOIN `groups` ON `groups`.id = trip_plans.group_id").
		Joins("JOIN user_groups ON user_groups.group_id = groups.id AND user_groups.user_id = ?", firebaseUID).
		Joins("JOIN cities ON cities.id = trip_plans.city_id").
		Where("trip_plans.deleted_at IS NULL").
		Scan(&tripPlans).Error
	if err != nil {
		return nil, err
	}

	return tripPlans, nil
}

func (r *TripPlanRepository) GetTripsByCityAndUser(cityId uint, firebaseUID string) ([]*responses.TripPlanResponse, error) {
	var tripPlans []*responses.TripPlanResponse

	err := r.DB.
		Table("trip_plans").
		Select("trip_plans.id AS id, trip_plans.name AS trip_name, trip_plans.start_date AS start_date, trip_plans.end_date AS end_date, groups.name AS group_name, cities.name AS city_name, cities.id as city_id").
		Joins("JOIN `groups` ON `groups`.id = trip_plans.group_id").
		Joins("JOIN user_groups ON user_groups.group_id = groups.id AND user_groups.user_id = ?", firebaseUID).
		Joins("JOIN cities ON cities.id = trip_plans.city_id").
		Where("trip_plans.deleted_at IS NULL AND cities.id = ?", cityId).
		Scan(&tripPlans).Error

	if err != nil {
		return nil, err
	}

	return tripPlans, nil
}

func (r *TripPlanRepository) GetTripsByGroupId(groupId uint) ([]*responses.TripPlanResponse, error) {
	var tripPlans []*responses.TripPlanResponse

	err := r.DB.
		Table("trip_plans").
		Select("trip_plans.id AS id, trip_plans.name AS trip_name, trip_plans.start_date AS start_date, trip_plans.end_date AS end_date, groups.name AS group_name, cities.name AS city_name, cities.id as city_id").
		Joins("JOIN `groups` ON `groups`.id = trip_plans.group_id").
		Joins("JOIN cities ON cities.id = trip_plans.city_id").
		Where("group_id = ? AND trip_plans.deleted_at IS NULL", groupId).
		Scan(&tripPlans).Error

	if err != nil {
		return nil, err
	}

	return tripPlans, nil
}
