package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type UserGroupController struct {
	UserGroupService services.UserGroupServiceInterface
}

func NewUserGroupController(userGroupService services.UserGroupServiceInterface) *UserGroupController {
	return &UserGroupController{
		UserGroupService: userGroupService,
	}
}

func (ugc *UserGroupController) GetAllByUserID(ginCtx *gin.Context) {
	userGroups, err := ugc.UserGroupService.GetGroupsByUserId(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginCtx.JSON(http.StatusOK, gin.H{"user_groups": userGroups})
}

func (ugc *UserGroupController) AddUserGroup(ginCtx *gin.Context) {
	var modifyUserGroupsReq requests.ModifyUserGroupRequest

	if err := ginCtx.ShouldBindJSON(&modifyUserGroupsReq); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	err := ugc.UserGroupService.AddUserGroup(ginCtx, &modifyUserGroupsReq)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{"message": "User group updated successfully"})
}

func (ugc *UserGroupController) DeleteUserGroup(ginCtx *gin.Context) {
	var modifyUserGroupsReq requests.ModifyUserGroupRequest

	if err := ginCtx.ShouldBindJSON(&modifyUserGroupsReq); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	err := ugc.UserGroupService.DeleteUserGroup(ginCtx, &modifyUserGroupsReq)
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{"message": "User group updated successfully"})
}

func (ugc *UserGroupController) GetAllUsersByGroupId(ctx *gin.Context) {
	groupIdString := ctx.Query("id")
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil || groupId < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}

	users, err := ugc.UserGroupService.GetAllUsersByGroupId(ctx, uint(groupId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}
