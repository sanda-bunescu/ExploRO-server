package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/common"
	"github.com/sanda-bunescu/ExploRO/initializers"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type UserServiceInterface interface {
	FindUserByEmail(email string) (*models.Users, error)
	RegisterUser(ctx context.Context, ginCtx *gin.Context) (*models.Users, error)
	LoginUser(ginCtx *gin.Context) (*models.Users, error)
	UpdateDeletedUser(user *models.Users, firebaseUID string) error
	SoftDelete(ctx context.Context) (*models.Users, error)
	GetUserByID(ctx context.Context, id any) (*models.Users, error)
	GetUserCitiesByUserID(ginCtx *gin.Context) ([]*responses.UserCityWithDetails, error)
	SoftDeleteUserGroupForUser(firebaseUID string) error
}

type UserService struct {
	UserRepo        repositories.UserRepositoryInterface
	UserCityService UserCityServiceInterface
	GroupService    GroupServiceInterface
	FirebaseService FirebaseServiceInterface
	ExpenseService  ExpenseServiceInterface
}

func NewUserService(repo repositories.UserRepositoryInterface, userCityService UserCityServiceInterface, groupService GroupServiceInterface, firebaseService FirebaseServiceInterface, ExpenseService ExpenseServiceInterface) *UserService {
	return &UserService{UserRepo: repo, UserCityService: userCityService, GroupService: groupService, FirebaseService: firebaseService, ExpenseService: ExpenseService}
}

var _ UserServiceInterface = (*UserService)(nil)

func (us *UserService) GetUserByID(ctx context.Context, id any) (*models.Users, error) {
	return us.UserRepo.GetByID(ctx, id)
}

func (us *UserService) FindUserByEmail(email string) (*models.Users, error) {
	user, err := us.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("UserService.FindUserByEmail: %v", err)
	}
	return user, nil
}

func (us *UserService) RegisterUser(ctx context.Context, ginCtx *gin.Context) (*models.Users, error) {
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user: no firebaseUID in context")
	}
	// Get userRecord from Firebase
	userRecord, err := us.FirebaseService.GetUserByUID(firebaseUID.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Firebase user details: %v", err)
	}
	// Check if user already exists
	user, err := us.UserRepo.GetUserByEmail(userRecord.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user from Database: %v", err)
	}

	if user != nil && user.DeletedAt == nil {
		return nil, fmt.Errorf("user already exists. Please log in instead")
	}

	// Create new user
	user = &models.Users{
		Id:    firebaseUID.(string),
		Name:  userRecord.DisplayName, //bodyWithUsername.Username
		Email: userRecord.Email,
	}
	if err := us.UserRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	//assign default cities
	defaultCities, err := common.SeedUserCities(initializers.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch default cities: %v", err)
	}

	if len(defaultCities) > 0 {
		for _, city := range defaultCities {
			defaultCity := models.UserCities{
				UserID: firebaseUID.(string),
				CityID: city.Id,
			}
			if err := us.UserRepo.AddDefaultCity(&defaultCity); err != nil {
				return nil, fmt.Errorf("failed to assign city %v to user: %v", city.Name, err)
			}
		}
	}
	//assign alone group
	newGroup := requests.NewGroup{
		Name: "Me",
	}

	err = us.GroupService.CreateGroup(ginCtx, newGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	return user, nil
}

func (us *UserService) LoginUser(ginCtx *gin.Context) (*models.Users, error) {
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user: no firebaseUID in context")
	}

	// Get userRecord from Firebase
	userRecord, err := us.FirebaseService.GetUserByUID(firebaseUID.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Firebase user details: %v", err)
	}

	// Check if user exists
	user, err := us.UserRepo.GetUserByEmail(userRecord.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user from Database: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user does not exist. Please register first")
	}

	return user, nil
}

func (us *UserService) UpdateDeletedUser(user *models.Users, firebaseUID string) error {
	if err := us.UserRepo.UpdateDeletedUser(user, firebaseUID); err != nil {
		return fmt.Errorf("UserService.UpdateDeletedUser: %v", err)
	}
	return nil
}

func (us *UserService) SoftDelete(ctx context.Context) (*models.Users, error) {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("failed to convert context to gin.Context")
	}
	user, err := us.FirebaseService.GetUserByFirebaseId(ginCtx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user does not exist")
	}

	//verify if account is deleted or not
	isDeleted, err := us.UserRepo.IsUserDeleted(user.Id)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, fmt.Errorf("user account is deleted")
	}

	err = us.UserCityService.DeleteUserCitiesByUserId(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user cities: %w", err)
	}

	err = us.ExpenseService.SoftDeleteExpenseByPayerId(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to soft delete expenses: %w", err)
	}

	err = us.GroupService.SoftDeleteGroupsWithOneActiveUser(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to soft delete groups: %w", err)
	}

	err = us.SoftDeleteUserGroupForUser(user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to soft delete user groups: %w", err)
	}

	err = us.UserRepo.SoftDelete(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("UserService.SoftDelete: %v", err)
	}
	return user, nil
}

func (us *UserService) GetUserCitiesByUserID(ginCtx *gin.Context) ([]*responses.UserCityWithDetails, error) {
	//get UID
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user: no firebaseUID in context")
	}
	//verify if account is deleted or not
	isDeleted, err := us.UserRepo.IsUserDeleted(firebaseUID.(string))
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, fmt.Errorf("user account is deleted")
	}
	//get cities
	existingUserCities, err := us.UserRepo.GetCitiesByUserId(firebaseUID.(string))
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUserCitiesByUserID: GetCitiesByUserId failed for user %s: %w", firebaseUID, err)
	}
	if len(existingUserCities) == 0 {
		return []*responses.UserCityWithDetails{}, nil
	}
	return existingUserCities, nil
}

func (us *UserService) SoftDeleteUserGroupForUser(firebaseUID string) error {
	err := us.UserRepo.SoftDeleteUserGroupForUser(firebaseUID)
	if err != nil {
		return fmt.Errorf("UserGroupService.SoftDeleteUserGroupForUser: %v", err)
	}

	return nil
}
