package services

import (
	"context"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type StopPointServiceInterface interface {
	GetAllTouristicAttractionsByItineraryId(itineraryId uint) ([]responses.StopPointResponse, error)
	AddTouristicAttractionsToItinerary(ctx context.Context, itineraryId uint, attractionIdsToAdd []uint) error
	Delete(ctx context.Context, stopPointId uint) error
}

type StopPointService struct {
	stopPointRepository repositories.StopPointRepositoryInterface
}

func NewStopPointService(stopPointRepository repositories.StopPointRepositoryInterface) *StopPointService {
	return &StopPointService{stopPointRepository: stopPointRepository}
}

var _ StopPointServiceInterface = (*StopPointService)(nil)

func (s *StopPointService) GetAllTouristicAttractionsByItineraryId(itineraryId uint) ([]responses.StopPointResponse, error) {
	var response []responses.StopPointResponse
	stopPoints, err := s.stopPointRepository.GetAllByItineraryId(itineraryId)
	if err != nil {
		return nil, err
	}

	for _, stopPoint := range stopPoints {
		touristicAttractionResponse := responses.TouristicAttractionResponse{
			TouristicAttractionId: stopPoint.TouristicAttraction.Id,
			Name:                  stopPoint.TouristicAttraction.Name,
			Description:           stopPoint.TouristicAttraction.Description,
			Category:              stopPoint.TouristicAttraction.Category,
			ImageUrl:              stopPoint.TouristicAttraction.ImageUrl,
			OpenHours:             stopPoint.TouristicAttraction.OpenHours,
			Fee:                   stopPoint.TouristicAttraction.Fee,
			Link:                  stopPoint.TouristicAttraction.Link,
		}
		stopPointResponse := responses.StopPointResponse{
			Id:                  stopPoint.Id,
			ItineraryId:         stopPoint.ItineraryId,
			TouristicAttraction: touristicAttractionResponse,
		}
		response = append(response, stopPointResponse)
	}
	return response, nil
}

func (s *StopPointService) AddTouristicAttractionsToItinerary(ctx context.Context, itineraryId uint, attractionIdsToAdd []uint) error {
	stopPoints, err := s.stopPointRepository.GetAllByItineraryId(itineraryId)
	if err != nil {
		return err
	}

	var nextStopPointNr int
	if len(stopPoints) > 0 {
		nextStopPointNr = stopPoints[len(stopPoints)-1].ItineraryStopNr + 1
	} else {
		nextStopPointNr = 1
	}

	for _, attractionId := range attractionIdsToAdd {
		stopPointToCreate := models.StopPoint{
			ItineraryStopNr: nextStopPointNr,
			ItineraryId:     itineraryId,
			AttractionId:    attractionId,
		}

		err = s.stopPointRepository.Create(ctx, &stopPointToCreate)
		if err != nil {
			return err
		}
		nextStopPointNr = nextStopPointNr + 1
	}
	return nil
}

func (s *StopPointService) Delete(ctx context.Context, stopPointId uint) error {
	stopPoint, err := s.stopPointRepository.GetByID(ctx, stopPointId)
	if err != nil {
		return err
	}

	err = s.stopPointRepository.SoftDelete(ctx, stopPoint)
	if err != nil {
		return err
	}
	return nil
}
