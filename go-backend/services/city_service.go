package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type CityServiceInterface interface {
	GetAllCities(ginCtx *gin.Context) ([]*responses.CityResponse, error)
}

type CityService struct {
	cityRepository repositories.CityRepositoryInterface
}

func NewCityService(cityRepo repositories.CityRepositoryInterface) *CityService {
	return &CityService{cityRepository: cityRepo}
}

var _ CityServiceInterface = (*CityService)(nil)

func (u *CityService) GetAllCities(ginCtx *gin.Context) ([]*responses.CityResponse, error) {
	cities, err := u.cityRepository.GetAll(ginCtx)
	if err != nil {
		return nil, fmt.Errorf("CityService.GetAllCities: GetAll failed : %w", err)
	}

	var cityResponses []*responses.CityResponse
	for _, city := range cities {
		cityResponse := &responses.CityResponse{
			CityId:          city.Id,
			CityName:        city.Name,
			CityDescription: city.Description,
			ImageURL:        city.ImageURL,
		}
		cityResponses = append(cityResponses, cityResponse)
	}

	return cityResponses, nil
}
