package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type StopPointController struct {
	StopPointService services.StopPointServiceInterface
}

func NewStopPointController(stopPointService services.StopPointServiceInterface) *StopPointController {
	return &StopPointController{StopPointService: stopPointService}
}

func (sc *StopPointController) GetAllByItineraryId(ctx *gin.Context) {
	var itineraryIdStr = ctx.Query("itineraryId")
	itineraryId, err := strconv.Atoi(itineraryIdStr)
	if err != nil || itineraryId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	stopPointsResponse, err := sc.StopPointService.GetAllTouristicAttractionsByItineraryId(uint(itineraryId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"stopPoints": stopPointsResponse})
}

func (sc *StopPointController) Create(ctx *gin.Context) {
	fmt.Println("Entered")
	var itineraryIdStr = ctx.Query("itineraryId")
	itineraryId, err := strconv.Atoi(itineraryIdStr)
	if err != nil || itineraryId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}

	var attractionIds *requests.TouristicAttractionIdsRequest
	err = ctx.ShouldBind(&attractionIds)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = sc.StopPointService.AddTouristicAttractionsToItinerary(ctx, uint(itineraryId), attractionIds.Ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "StopPoints created successfully"})
}

func (sc *StopPointController) Delete(ctx *gin.Context) {
	var stopPointIdStr = ctx.Query("stopPointId")
	stopPointId, err := strconv.Atoi(stopPointIdStr)
	if err != nil || stopPointId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	err = sc.StopPointService.Delete(ctx, uint(stopPointId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "StopPoints deleted successfully"})

}
