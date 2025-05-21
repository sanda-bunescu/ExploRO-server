package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type UserController struct {
	UserService     services.UserServiceInterface
	FirebaseService services.FirebaseServiceInterface
	UserCityService services.UserCityServiceInterface
}

func NewUserController(userService services.UserServiceInterface, firebaseService services.FirebaseServiceInterface, userCityService services.UserCityServiceInterface) *UserController {
	return &UserController{
		UserService:     userService,
		FirebaseService: firebaseService,
		UserCityService: userCityService,
	}
}

func (uc *UserController) Register(ctx *gin.Context) {
	user, err := uc.UserService.RegisterUser(context.Background(), ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": user})
}

func (uc *UserController) Login(ctx *gin.Context) {
	user, err := uc.UserService.LoginUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "user": user})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {

	user, err := uc.UserService.SoftDelete(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted", "user": user})
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	idParam := ctx.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID" + err.Error()})
		return
	}
	user, err := uc.UserService.GetUserByID(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *UserController) GetUserCitiesByUserID(ginCtx *gin.Context) {
	// Call the service layer
	userCities, err := uc.UserService.GetUserCitiesByUserID(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the fetched cities
	ginCtx.JSON(http.StatusOK, gin.H{"user_cities": userCities})
}

func (uc *UserController) AddUserCity(ctx *gin.Context) {
	var req requests.CityRequest
	//get data from JSON/body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data provided", "details": err.Error()})
		return
	}
	err := uc.UserCityService.AddUserCity(ctx, req.CityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "City added to user successfully"})
}

func (uc *UserController) DeleteUserCity(ctx *gin.Context) {
	var req requests.CityRequest
	//get data from JSON/body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data provided", "details": err.Error()})
		return
	}

	err := uc.UserCityService.DeleteUserCity(ctx, req.CityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "City removed from user successfully"})
}
