package repositories

import (
	"fmt"
	"github.com/sanda-bunescu/ExploRO/models"
	"gorm.io/gorm"
	"time"
)

type GroupRepositoryInterface interface {
	SoftDeleteGroupsWithOneActiveUser(firebaseUID string) ([]uint, error)
	SoftDeleteGroupWhenAloneUserLeaves(groupId uint, firebaseUID string) error
	BaseRepositoryInterface[models.Group]
}

type GroupRepository struct {
	BaseRepository[models.Group]
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{BaseRepository[models.Group]{DB: db}}
}

var _ GroupRepositoryInterface = (*GroupRepository)(nil)

func (gr *GroupRepository) SoftDeleteGroupsWithOneActiveUser(firebaseUID string) ([]uint, error) {
	now := time.Now()
	var groupIDs []uint
	err := gr.DB.Table("user_groups ug").
		Where("ug.user_id = ? AND ug.deleted_at IS NULL", firebaseUID).
		Where(`
			(SELECT COUNT(*) 
			FROM user_groups 
			WHERE group_id = ug.group_id 
			AND deleted_at IS NULL) = 1
		`).
		Select("group_id").
		Find(&groupIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch group IDs: %w", err)
	}

	if len(groupIDs) == 0 {
		return []uint{}, nil
	}

	err = gr.DB.Model(&models.Group{}).Where("id IN ?", groupIDs).Update("deleted_at", now).Error
	if err != nil {
		return nil, fmt.Errorf("failed to soft delete groups: %w", err)
	}

	return groupIDs, nil
}

func (gr *GroupRepository) SoftDeleteGroupWhenAloneUserLeaves(groupId uint, firebaseUID string) error {
	now := time.Now()

	err := gr.DB.Exec(`
		UPDATE `+"`groups`"+`
		SET deleted_at = ?
		WHERE id IN (
			SELECT ug.group_id
		FROM user_groups ug
		WHERE ug.user_id = ?
		AND ug.group_id = ?
		AND ug.deleted_at IS NULL
		AND (
				SELECT COUNT(*)
				FROM user_groups
				WHERE group_id = ug.group_id
				AND deleted_at IS NULL
			) = 1
		);
	`, now, firebaseUID, groupId).Error

	if err != nil {
		return fmt.Errorf("failed to soft delete groups with one active user (%s): %w", firebaseUID, err)
	}

	return nil
}
