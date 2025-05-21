package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/repositories"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type GroupServiceInterface interface {
	CreateGroup(ginCtx *gin.Context, group requests.NewGroup) error
	SoftDelete(ctx context.Context, groupId uint) error
	SoftDeleteGroupsWithOneActiveUser(ctx context.Context, userId string) error
	SoftDeleteGroupWhenAloneUserLeavesGroup(ctx context.Context, groupId uint, userId string) error
	GetById(ctx context.Context, groupId uint) (*models.Group, error)
}

type GroupService struct {
	groupRepository     repositories.GroupRepositoryInterface
	FirebaseService     FirebaseServiceInterface
	TripPlanService     TripPlanServiceInterface
	userGroupRepository repositories.UserGroupRepositoryInterface
	expenseService      ExpenseServiceInterface
}

func NewGroupService(groupRepository repositories.GroupRepositoryInterface, firebaseService FirebaseServiceInterface, tripPlanService TripPlanServiceInterface, userGroupRepository repositories.UserGroupRepositoryInterface, expenseService ExpenseServiceInterface) *GroupService {
	return &GroupService{groupRepository: groupRepository, FirebaseService: firebaseService, TripPlanService: tripPlanService, userGroupRepository: userGroupRepository, expenseService: expenseService}
}

var _ GroupServiceInterface = (*GroupService)(nil)

func (g *GroupService) GetById(ctx context.Context, groupId uint) (*models.Group, error) {
	group, err := g.groupRepository.GetByID(ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("group: %w", err)
	}
	return group, nil
}

func (g *GroupService) CreateGroup(ginCtx *gin.Context, newGroup requests.NewGroup) error {
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return fmt.Errorf("unauthorized user: no firebaseUID in context")
	}

	imageURL, err := getRandomGroupImage()
	if err != nil {
		return fmt.Errorf("could not assign random image: %w", err)
	}

	group := &models.Group{
		Name:     newGroup.Name,
		ImageURL: imageURL,
	}
	err = g.groupRepository.Create(ginCtx, group)
	if err != nil {
		return err
	}

	//add user to group
	newUserGroup := &models.UserGroup{
		UserId:  firebaseUID.(string),
		GroupId: group.Id,
	}
	err = g.userGroupRepository.Create(ginCtx, newUserGroup)
	if err != nil {
		return fmt.Errorf("GroupService.CreateGroup: Failed to add user to group, Error: %w", err)
	}
	return nil
}

func (g *GroupService) SoftDelete(ctx context.Context, groupId uint) error {
	group, err := g.groupRepository.GetByID(ctx, groupId)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDelete: %v", err)
	}

	// Get all user-group relations for this group
	userGroups, err := g.userGroupRepository.GetByGroupId(groupId)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDelete: failed to retrieve user-group relations: %w", err)
	}
	//soft delete userGroups
	for _, userGroup := range userGroups {
		err = g.userGroupRepository.SoftDelete(ctx, userGroup)
		if err != nil {
			return fmt.Errorf("GroupService.SoftDelete: failed to delete user-group relation for user %v: %w", userGroup.UserId, err)
		}
	}

	// Get all trip plans for this group
	tripPlans, err := g.TripPlanService.GetTripsByGroupId(groupId)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDelete: failed to retrieve trip plans: %w", err)
	}

	// Soft delete trip plans
	for _, tripPlan := range tripPlans {
		err = g.TripPlanService.DeleteTrip(ctx, tripPlan.Id)
		if err != nil {
			return fmt.Errorf("GroupService.SoftDelete: failed to delete trip plan %v: %w", tripPlan.Id, err)
		}
	}
	// Get all expensesfor this group
	expenses, err := g.expenseService.GetAllExpensesByGroupId(ctx, groupId)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDelete: failed to retrieve expenses for group %v: %w", groupId, err)
	}
	for _, expense := range expenses {
		if err := g.expenseService.SoftDeleteExpenseByID(ctx, expense.Id); err != nil {
			return fmt.Errorf("GroupService.SoftDelete: failed to delete expense %v: %w", expense.Id, err)
		}
	}
	//soft delete group
	err = g.groupRepository.SoftDelete(ctx, group)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDelete: %v", err)
	}

	return nil
}

func (g *GroupService) SoftDeleteGroupsWithOneActiveUser(ctx context.Context, userId string) error {
	groupIds, err := g.groupRepository.SoftDeleteGroupsWithOneActiveUser(userId)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDeleteGroupsWithOneActiveUser: %w", err)
	}
	for _, groupId := range groupIds {
		// Get all trip plans for this group
		tripPlans, err := g.TripPlanService.GetTripsByGroupId(groupId)
		if err != nil {
			return fmt.Errorf("GroupService.SoftDelete: failed to retrieve trip plans for group %v: %w", groupId, err)
		}

		// Soft delete trip plans
		for _, tripPlan := range tripPlans {
			err = g.TripPlanService.DeleteTrip(ctx, tripPlan.Id)
			if err != nil {
				return fmt.Errorf("GroupService.SoftDelete: failed to delete trip plan %v for group %v: %w", tripPlan.Id, groupId, err)
			}
		}

	}
	return nil
}

func (g *GroupService) SoftDeleteGroupWhenAloneUserLeavesGroup(ctx context.Context, groupId uint, userId string) error {
	// Get all trip plans for this group
	tripPlans, err := g.TripPlanService.GetTripsByGroupId(groupId)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDelete: failed to retrieve trip plans: %w", err)
	}

	// Soft delete trip plans
	for _, tripPlan := range tripPlans {
		err = g.TripPlanService.DeleteTrip(ctx, tripPlan.Id)
		if err != nil {
			return fmt.Errorf("GroupService.SoftDelete: failed to delete trip plan %v: %w", tripPlan.Id, err)
		}
	}
	// Get all expensesfor this group
	expenses, err := g.expenseService.GetAllExpensesByGroupId(ctx, groupId)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDelete: failed to retrieve expenses for group %v: %w", groupId, err)
	}
	for _, expense := range expenses {
		if err := g.expenseService.SoftDeleteExpenseByID(ctx, expense.Id); err != nil {
			return fmt.Errorf("GroupService.SoftDelete: failed to delete expense %v: %w", expense.Id, err)
		}
	}

	err = g.groupRepository.SoftDeleteGroupWhenAloneUserLeaves(groupId, userId)
	if err != nil {
		return fmt.Errorf("GroupService.SoftDeleteGroupsWithOneActiveUser: %w", err)
	}
	return nil
}

func getRandomGroupImage() (string, error) {

	files, err := os.ReadDir("static/groupImages")
	if err != nil {
		return "", err
	}

	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".jpg" || filepath.Ext(file.Name()) == ".png" || filepath.Ext(file.Name()) == ".jpeg") {
			imageFiles = append(imageFiles, file.Name())
		}
	}

	if len(imageFiles) == 0 {
		return "", fmt.Errorf("no images found in groupImages")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	chosen := imageFiles[rng.Intn(len(imageFiles))]

	return "/static/groupImages/" + chosen, nil
}
