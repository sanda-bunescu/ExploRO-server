package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
)

type CityController struct {
	CityService services.CityServiceInterface
}

func NewCityController(cityService services.CityServiceInterface) *CityController {
	return &CityController{
		CityService: cityService,
	}
}

func (c *CityController) GetAllCities(ctx *gin.Context) {
	cities, err := c.CityService.GetAllCities(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch cities: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"cities": cities})
}
