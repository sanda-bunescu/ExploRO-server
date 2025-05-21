package services

import (
	"context"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type ItineraryServiceInterface interface {
	GetAllItinerariesByTripPlanId(tripPlanId uint) ([]responses.ItineraryResponse, error)
	CreateItineraryInTripPlan(ctx context.Context, tripPlanId uint) (*responses.ItineraryResponse, error)
	DeleteItinerary(ctx context.Context, itineraryId uint) error
}

type ItineraryService struct {
	itineraryRepository repositories.ItineraryRepositoryInterface
	stopPointService    StopPointServiceInterface
}

func NewItineraryService(itineraryRepository repositories.ItineraryRepositoryInterface, stopPointService StopPointServiceInterface) *ItineraryService {
	return &ItineraryService{itineraryRepository: itineraryRepository, stopPointService: stopPointService}
}

var _ ItineraryServiceInterface = (*ItineraryService)(nil)

func (is *ItineraryService) GetAllItinerariesByTripPlanId(tripPlanId uint) ([]responses.ItineraryResponse, error) {
	var itinerariesResponse []responses.ItineraryResponse
	itineraries, err := is.itineraryRepository.GetAllByTripPlanId(tripPlanId)
	if err != nil {
		return nil, err
	}
	for _, itinerary := range itineraries {
		itineraryResponse := responses.ItineraryResponse{
			Id:    itinerary.Id,
			DayNr: itinerary.DayNr,
		}
		itinerariesResponse = append(itinerariesResponse, itineraryResponse)
	}
	return itinerariesResponse, nil
}

func (is *ItineraryService) CreateItineraryInTripPlan(ctx context.Context, tripPlanId uint) (*responses.ItineraryResponse, error) {
	var itineraryResponse *responses.ItineraryResponse
	// Fetch itineraries by tripPlanId and order by Name ("Day X" format)
	lastAddedItinerary, err := is.itineraryRepository.GetLastAddedItinerary(tripPlanId)
	if err != nil {
		return nil, err
	}
	// If no itinerary exists, create "Day 1"
	if lastAddedItinerary == nil {
		newItinerary := models.Itinerary{
			DayNr:      1,
			TripPlanId: tripPlanId,
		}

		err = is.itineraryRepository.Create(ctx, &newItinerary)
		if err != nil {
			return nil, err
		}

		itineraryResponse = &responses.ItineraryResponse{
			Id:    newItinerary.Id,
			DayNr: newItinerary.DayNr,
		}

		return itineraryResponse, nil
	} else {
		dayNum := lastAddedItinerary.DayNr + 1

		newItinerary := models.Itinerary{
			DayNr:      dayNum,
			TripPlanId: tripPlanId,
		}

		err = is.itineraryRepository.Create(ctx, &newItinerary)
		if err != nil {
			return nil, err
		}
		itineraryResponse = &responses.ItineraryResponse{
			Id:    newItinerary.Id,
			DayNr: newItinerary.DayNr,
		}
		return itineraryResponse, nil
	}
}

func (is *ItineraryService) DeleteItinerary(ctx context.Context, itineraryId uint) error {
	itinerary, err := is.itineraryRepository.GetByID(ctx, itineraryId)
	if err != nil {
		return err
	}

	stopPoints, err := is.stopPointService.GetAllTouristicAttractionsByItineraryId(itineraryId)
	if err != nil {
		return err
	}

	for _, stopPoint := range stopPoints {
		err := is.stopPointService.Delete(ctx, stopPoint.Id)
		if err != nil {
			return err
		}
	}

	err = is.itineraryRepository.SoftDelete(ctx, itinerary)
	if err != nil {
		return err
	}
	return nil
}
