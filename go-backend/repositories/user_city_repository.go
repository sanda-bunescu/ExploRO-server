package repositories

import (
	"fmt"
	"github.com/sanda-bunescu/ExploRO/models"
	"gorm.io/gorm"
)

type UserCityRepositoryInterface interface {
	GetUserCity(userID string, cityID uint) (*models.UserCities, error)
	GetUserCitiesByUserID(firebaseUID string) ([]*models.UserCities, error)
	BaseRepositoryInterface[models.UserCities]
}

type UserCityRepository struct {
	BaseRepository[models.UserCities]
}

func NewUserCityRepository(db *gorm.DB) *UserCityRepository {
	return &UserCityRepository{BaseRepository: BaseRepository[models.UserCities]{DB: db}}
}

var _ UserCityRepositoryInterface = (*UserCityRepository)(nil)

func (ur *UserCityRepository) GetUserCity(userID string, cityID uint) (*models.UserCities, error) {
	var userCity models.UserCities
	err := ur.DB.Where("user_id = ? AND city_id = ? AND deleted_at IS NULL", userID, cityID).First(&userCity).Error
	if err != nil {
		return nil, err
	}
	return &userCity, nil
}

func (ur *UserCityRepository) GetUserCitiesByUserID(firebaseUID string) ([]*models.UserCities, error) {
	var userCities []*models.UserCities

	// Fetch all user-city relationships for the given user
	err := ur.DB.Where("user_id = ? AND deleted_at IS NULL", firebaseUID).Find(&userCities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user cities: %w", err)
	}

	return userCities, nil
}
