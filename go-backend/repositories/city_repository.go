package repositories

import (
	"github.com/sanda-bunescu/ExploRO/models"
	"gorm.io/gorm"
)

type CityRepositoryInterface interface {
	GetCityById(id uint) (*models.City, error)
	BaseRepositoryInterface[models.City]
}

type CityRepository struct {
	BaseRepository[models.City]
}

func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{BaseRepository: BaseRepository[models.City]{DB: db}}
}

var _ CityRepositoryInterface = (*CityRepository)(nil)

func (c *CityRepository) GetCityById(id uint) (*models.City, error) {
	var city *models.City
	err := c.DB.Where("id = ?", id).First(&city).Error
	if err != nil {
		return nil, err
	}
	return city, nil
}
