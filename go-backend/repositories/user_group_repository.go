package repositories

import (
	"fmt"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"gorm.io/gorm"
)

type UserGroupRepositoryInterface interface {
	GetExistingGroupUsersIds(groupId uint) ([]string, error)
	GetGroupsByUserId(userId string) ([]*responses.UserGroupResponse, error)
	GetUserGroupByUserAndGroup(userId string, groupId uint) (*models.UserGroup, error)
	GetByGroupId(groupId uint) ([]*models.UserGroup, error)
	GetAllUsersByGroupId(groupId uint) ([]*responses.GroupUserResponse, error)
	IsUserInGroup(userId string, groupId uint) (bool, error)
	BaseRepositoryInterface[models.UserGroup]
}

type UserGroupRepository struct {
	BaseRepository[models.UserGroup]
}

func NewUserGroupRepository(db *gorm.DB) *UserGroupRepository {
	return &UserGroupRepository{BaseRepository[models.UserGroup]{DB: db}}
}

var _ UserGroupRepositoryInterface = (*UserGroupRepository)(nil)

func (ug *UserGroupRepository) GetGroupsByUserId(userId string) ([]*responses.UserGroupResponse, error) {
	var userGroupsResponse []*responses.UserGroupResponse
	err := ug.DB.Table("user_groups").
		Select("`groups`.id AS id, `groups`.name AS group_name, `groups`.image_url AS image_url").
		Joins("INNER JOIN `groups` ON `groups`.id = user_groups.group_id").
		Where("user_groups.user_id = ? AND user_groups.deleted_at IS NULL ", userId).
		Find(&userGroupsResponse).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user cities for userID %s: %w", userId, err)
	}
	return userGroupsResponse, nil
}

func (ug *UserGroupRepository) GetExistingGroupUsersIds(groupId uint) ([]string, error) {
	var userGroupsIds []string

	err := ug.DB.Table("users").
		Select("users.id").
		Joins("INNER JOIN user_groups ON users.id = user_groups.user_id").
		Where("user_groups.group_id = ? AND user_groups.deleted_at IS NULL", groupId).
		Find(&userGroupsIds).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch users for group ID %d: %w", groupId, err)
	}
	return userGroupsIds, nil
}

func (ug *UserGroupRepository) GetUserGroupByUserAndGroup(userId string, groupId uint) (*models.UserGroup, error) {
	var userGroup models.UserGroup

	err := ug.DB.
		Where("user_id = ? AND group_id = ? AND deleted_at IS NULL", userId, groupId).
		First(&userGroup).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user group for userID %s and groupID %d: %w", userId, groupId, err)
	}

	return &userGroup, nil
}

func (ug *UserGroupRepository) GetByGroupId(groupId uint) ([]*models.UserGroup, error) {
	var userGroups []*models.UserGroup

	err := ug.DB.
		Where("group_id = ? AND deleted_at IS NULL", groupId).
		Find(&userGroups).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user groups for groupID %d: %w", groupId, err)
	}

	return userGroups, nil
}

func (ug *UserGroupRepository) GetAllUsersByGroupId(groupId uint) ([]*responses.GroupUserResponse, error) {
	var users []*responses.GroupUserResponse

	err := ug.DB.
		Table("user_groups").
		Select("users.id as user_id, users.name as user_name, users.email as user_email").
		Joins("JOIN users ON users.id = user_groups.user_id").
		Where("user_groups.group_id = ? AND user_groups.deleted_at IS NULL", groupId).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch users for groupID %d: %w", groupId, err)
	}
	return users, nil
}

func (ug *UserGroupRepository) IsUserInGroup(userId string, groupId uint) (bool, error) {
	var count int64

	err := ug.DB.
		Table("user_groups").
		Where("user_id = ? AND group_id = ? AND deleted_at IS NULL", userId, groupId).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check if user is in group: %w", err)
	}

	return count > 0, nil
}
