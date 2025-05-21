package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type GroupController struct {
	GroupService services.GroupServiceInterface
}

func NewGroupController(groupService services.GroupServiceInterface) *GroupController {
	return &GroupController{
		GroupService: groupService,
	}
}

func (gc *GroupController) CreateGroup(ctx *gin.Context) {
	var req requests.NewGroup

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data provided", "details": err.Error()})
	}

	err := gc.GroupService.CreateGroup(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Group created successfully"})
}

func (gc *GroupController) DeleteGroup(ctx *gin.Context) {
	groupIdString := ctx.Query("id")
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil || groupId < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	err = gc.GroupService.SoftDelete(ctx, uint(groupId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}
