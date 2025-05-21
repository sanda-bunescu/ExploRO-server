package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type TouristicAttractionServiceInterface interface {
	GetAllByCityId(cityId uint) ([]*responses.TouristicAttractionResponse, error)
	GetAllTouristicAttractions(ginCtx *gin.Context) ([]*responses.TouristicAttractionResponse, error)
	GetAttractionsNotInItinerary(cityId uint, tripPlanId uint) ([]*responses.TouristicAttractionResponse, error)
}

type TouristicAttractionService struct {
	touristicAttractionsRepository repositories.TouristicAttractionRepositoryInterface
}

func NewTouristicAttractionService(TouristicAttractionRepo repositories.TouristicAttractionRepositoryInterface) *TouristicAttractionService {
	return &TouristicAttractionService{touristicAttractionsRepository: TouristicAttractionRepo}
}

var _ TouristicAttractionServiceInterface = (*TouristicAttractionService)(nil)

func (t TouristicAttractionService) GetAllTouristicAttractions(ginCtx *gin.Context) ([]*responses.TouristicAttractionResponse, error) {
	touristicAttractions, err := t.touristicAttractionsRepository.GetAll(ginCtx)
	if err != nil {
		return nil, fmt.Errorf("TouristicAttractionService.GetAllTouristicAttractions failed: %w", err)
	}

	var touristicAttractionsResponse []*responses.TouristicAttractionResponse
	for _, touristicAttraction := range touristicAttractions {
		touristicAttractionResponse := &responses.TouristicAttractionResponse{
			TouristicAttractionId: touristicAttraction.Id,
			Name:                  touristicAttraction.Name,
			Description:           touristicAttraction.Description,
			Category:              touristicAttraction.Category,
			ImageUrl:              touristicAttraction.ImageUrl,
			OpenHours:             touristicAttraction.OpenHours,
			Fee:                   touristicAttraction.Fee,
			Link:                  touristicAttraction.Link,
		}
		touristicAttractionsResponse = append(touristicAttractionsResponse, touristicAttractionResponse)
	}
	return touristicAttractionsResponse, nil
}

func (t TouristicAttractionService) GetAllByCityId(cityId uint) ([]*responses.TouristicAttractionResponse, error) {
	touristicAttractions, err := t.touristicAttractionsRepository.GetAttractionsByCityId(cityId)
	if err != nil {
		return nil, fmt.Errorf("TouristicAttractionService.GetAllByCityId failed: %w", err)
	}
	var touristicAttractionsResponse []*responses.TouristicAttractionResponse
	for _, touristicAttraction := range touristicAttractions {
		touristicAttractionResponse := &responses.TouristicAttractionResponse{
			TouristicAttractionId: touristicAttraction.Id,
			Name:                  touristicAttraction.Name,
			Description:           touristicAttraction.Description,
			Category:              touristicAttraction.Category,
			ImageUrl:              touristicAttraction.ImageUrl,
			OpenHours:             touristicAttraction.OpenHours,
			Fee:                   touristicAttraction.Fee,
			Link:                  touristicAttraction.Link,
		}
		touristicAttractionsResponse = append(touristicAttractionsResponse, touristicAttractionResponse)
	}
	return touristicAttractionsResponse, nil
}

func (t TouristicAttractionService) GetAttractionsNotInItinerary(cityId uint, tripPlanId uint) ([]*responses.TouristicAttractionResponse, error) {

	touristicAttractions, err := t.touristicAttractionsRepository.GetAttractionsThatAreNotInItineraryStopPoint(cityId, tripPlanId)
	if err != nil {
		return nil, fmt.Errorf("TouristicAttractionService.GetAttractionsNotInItinerary failed: %w", err)
	}

	var touristicAttractionsResponse []*responses.TouristicAttractionResponse
	for _, touristicAttraction := range touristicAttractions {
		touristicAttractionsResponse = append(touristicAttractionsResponse, &responses.TouristicAttractionResponse{
			TouristicAttractionId: touristicAttraction.Id,
			Name:                  touristicAttraction.Name,
			Description:           touristicAttraction.Description,
			Category:              touristicAttraction.Category,
			ImageUrl:              touristicAttraction.ImageUrl,
			OpenHours:             touristicAttraction.OpenHours,
			Fee:                   touristicAttraction.Fee,
			Link:                  touristicAttraction.Link,
		})
	}

	return touristicAttractionsResponse, nil
}
