package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type TouristicAttractionController struct {
	TouristicAttractionService services.TouristicAttractionServiceInterface
}

func NewTouristicAttractionController(service services.TouristicAttractionServiceInterface) *TouristicAttractionController {
	return &TouristicAttractionController{TouristicAttractionService: service}
}

func (t *TouristicAttractionController) GetAllTouristicAttractions(ginCtx *gin.Context) {
	touristicAttractions, err := t.TouristicAttractionService.GetAllTouristicAttractions(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, gin.H{"touristic_attractions": touristicAttractions})
}

func (t *TouristicAttractionController) GetAllByCityId(ctx *gin.Context) {
	cityIDString := ctx.Query("cityId")

	cityID, err := strconv.Atoi(cityIDString)
	if err != nil || cityID < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cityId parameter"})
		return
	}
	touristicAttractions, err := t.TouristicAttractionService.GetAllByCityId(uint(cityID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"touristic_attractions": touristicAttractions})
}

func (t *TouristicAttractionController) GetAttractionsNotInItinerary(ctx *gin.Context) {
	cityIDString := ctx.Query("cityId")
	cityID, err := strconv.Atoi(cityIDString)
	if err != nil || cityID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cityId parameter"})
		return
	}

	tripPlanIDString := ctx.Query("tripPlanId")
	tripPlanID, err := strconv.Atoi(tripPlanIDString)
	if err != nil || tripPlanID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tripPlanId parameter"})
		return
	}

	touristicAttractions, err := t.TouristicAttractionService.GetAttractionsNotInItinerary(uint(cityID), uint(tripPlanID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"touristic_attractions": touristicAttractions})
}
