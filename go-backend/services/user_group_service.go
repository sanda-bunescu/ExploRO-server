package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type UserGroupServiceInterface interface {
	GetGroupsByUserId(ginCtx *gin.Context) ([]*responses.UserGroupResponse, error)
	AddUserGroup(ginCtx *gin.Context, userGroup *requests.ModifyUserGroupRequest) error
	DeleteUserGroup(ginCtx *gin.Context, userGroup *requests.ModifyUserGroupRequest) error
	GetAllUsersByGroupId(ginCtx *gin.Context, groupId uint) ([]*responses.GroupUserResponse, error)
}

type UserGroupService struct {
	UserRepo            repositories.UserRepositoryInterface
	GroupService        GroupServiceInterface
	UserGroupRepository repositories.UserGroupRepositoryInterface
}

func NewUserGroupService(userGroupRepository repositories.UserGroupRepositoryInterface, GroupService GroupServiceInterface, UserRepo repositories.UserRepositoryInterface) *UserGroupService {
	return &UserGroupService{
		UserRepo:            UserRepo,
		GroupService:        GroupService,
		UserGroupRepository: userGroupRepository,
	}
}

var _ UserGroupServiceInterface = (*UserGroupService)(nil)

func (ugs *UserGroupService) GetGroupsByUserId(ginCtx *gin.Context) ([]*responses.UserGroupResponse, error) {
	//get UID
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user: no firebaseUID in context")
	}
	//verify if account is deleted or not
	isDeleted, err := ugs.UserRepo.IsUserDeleted(firebaseUID.(string))
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, fmt.Errorf("user account is deleted")
	}
	//get groups
	existingUserGroups, err := ugs.UserGroupRepository.GetGroupsByUserId(firebaseUID.(string))

	if err != nil {
		return nil, fmt.Errorf("UserGroupService.GetGroupsByUserId: GetGroupsByUserId failed for user %s: %w", firebaseUID, err)
	}
	if len(existingUserGroups) == 0 {
		return []*responses.UserGroupResponse{}, nil
	}
	return existingUserGroups, nil
}

func (ugs *UserGroupService) AddUserGroup(ginCtx *gin.Context, userGroup *requests.ModifyUserGroupRequest) error {
	userEmailToAdd := userGroup.UserEmail

	userId, err := ugs.UserRepo.GetUserIdByEmail(userEmailToAdd)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if userId == "" {
		return fmt.Errorf("no user found with email %s", userEmailToAdd)
	}

	existingUserGroup, err := ugs.UserGroupRepository.GetUserGroupByUserAndGroup(userId, userGroup.GroupId)
	if err == nil && existingUserGroup != nil {
		return fmt.Errorf("user is already in the group")
	}
	newUserGroup := &models.UserGroup{
		UserId:  userId,
		GroupId: userGroup.GroupId,
	}
	err = ugs.UserGroupRepository.Create(ginCtx, newUserGroup)
	if err != nil {
		return fmt.Errorf("failed to add user to group, Error: %w", err)
	}

	return nil
}

func (ugs *UserGroupService) DeleteUserGroup(ginCtx *gin.Context, userGroup *requests.ModifyUserGroupRequest) error {
	userEmailToAdd := userGroup.UserEmail

	userId, err := ugs.UserRepo.GetUserIdByEmail(userEmailToAdd)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	existingUserGroup, err := ugs.UserGroupRepository.GetUserGroupByUserAndGroup(userId, userGroup.GroupId)
	if err != nil {
		return fmt.Errorf("UserGroupService.UpdateUserGroup: GetUserGroupByUserAndGroup failed, Error: %w", err)
	}
	//if user is alone in group, at leave group action will delete group
	err = ugs.GroupService.SoftDeleteGroupWhenAloneUserLeavesGroup(ginCtx, userGroup.GroupId, userId)
	if err != nil {
		return fmt.Errorf("failed to soft delete groups: %w", err)
	}

	err = ugs.UserGroupRepository.SoftDelete(ginCtx, existingUserGroup)
	if err != nil {
		return fmt.Errorf("UserGroupService.UpdateUserGroup: SoftDelete failed, Error: %w", err)
	}

	return nil
}

func (ugs *UserGroupService) GetAllUsersByGroupId(ginCtx *gin.Context, groupId uint) ([]*responses.GroupUserResponse, error) {
	//get UID
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user: no firebaseUID in context")
	}
	//verify if account is deleted or not
	isDeleted, err := ugs.UserRepo.IsUserDeleted(firebaseUID.(string))
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, fmt.Errorf("user account is deleted")
	}

	isMember, err := ugs.UserGroupRepository.IsUserInGroup(firebaseUID.(string), groupId)
	if err != nil {
		return nil, fmt.Errorf("UserGroupService.GetAllUsersByGroupId: %w", err)
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of the group")
	}

	users, err := ugs.UserGroupRepository.GetAllUsersByGroupId(groupId)
	if err != nil {
		return nil, fmt.Errorf("UserGroupService.GetAllUsersByGroupId: %w", err)
	}
	return users, nil
}
