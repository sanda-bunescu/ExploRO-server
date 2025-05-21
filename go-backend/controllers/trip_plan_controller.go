package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type TripPlanController struct {
	TripPlanService services.TripPlanServiceInterface
}

func NewTripPlanController(tripPlanService services.TripPlanServiceInterface) *TripPlanController {
	return &TripPlanController{
		TripPlanService: tripPlanService,
	}
}

func (t *TripPlanController) GetTripsByUserId(ctx *gin.Context) {
	trips, err := t.TripPlanService.GetTripsByUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"trips": trips})
}

func (t *TripPlanController) GetTripsByCityAndUser(ctx *gin.Context) {
	cityIdString := ctx.Query("cityId")
	cityId, err := strconv.Atoi(cityIdString)
	if err != nil || cityId < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}

	trips, err := t.TripPlanService.GetTripsByCityAndUser(uint(cityId), ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"trips": trips})
}

func (t *TripPlanController) GetTripsByGroupId(ctx *gin.Context) {
	groupIdString := ctx.Query("groupId")
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil || groupId < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}

	trips, err := t.TripPlanService.GetTripsByGroupId(uint(groupId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"trips": trips})
}

func (t *TripPlanController) CreateTrip(ctx *gin.Context) {
	var tripRequest requests.CreateTripPlanRequest
	if err := ctx.ShouldBindJSON(&tripRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := t.TripPlanService.CreateTrip(context.Background(), &tripRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Trip created successfully"})
}

func (t *TripPlanController) DeleteTrip(ctx *gin.Context) {
	tripIdString := ctx.Query("tripId")
	tripId, err := strconv.Atoi(tripIdString)
	if err != nil || tripId < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}

	err = t.TripPlanService.DeleteTrip(context.Background(), uint(tripId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Trip deleted successfully"})
}
