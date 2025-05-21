package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type ItineraryController struct {
	ItineraryService services.ItineraryServiceInterface
}

func NewItineraryController(itineraryService services.ItineraryServiceInterface) *ItineraryController {
	return &ItineraryController{ItineraryService: itineraryService}
}

func (it *ItineraryController) GetByTripPlanId(ctx *gin.Context) {
	var planIdStr = ctx.Query("tripPlanId")
	planId, err := strconv.Atoi(planIdStr)
	if err != nil || planId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	itinerariesResponse, err := it.ItineraryService.GetAllItinerariesByTripPlanId(uint(planId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"itineraries": itinerariesResponse})
}

func (it *ItineraryController) Create(ctx *gin.Context) {
	var planIdStr = ctx.Query("tripPlanId")
	planId, err := strconv.Atoi(planIdStr)
	if err != nil || planId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	itineraryResponse, err := it.ItineraryService.CreateItineraryInTripPlan(ctx, uint(planId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"itinerary": itineraryResponse})
}

func (it *ItineraryController) Delete(ctx *gin.Context) {
	var itineraryIdStr = ctx.Query("itineraryId")
	itineraryId, err := strconv.Atoi(itineraryIdStr)
	if err != nil || itineraryId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	err = it.ItineraryService.DeleteItinerary(ctx, uint(itineraryId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Itinerary deleted successfully"})
}
