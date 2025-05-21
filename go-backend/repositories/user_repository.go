package repositories

import (
	"errors"
	"fmt"
	"github.com/sanda-bunescu/ExploRO/initializers"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"gorm.io/gorm"
	"time"
)

type UserRepositoryInterface interface {
	GetUserByEmail(email string) (*models.Users, error)
	GetUserByFirebaseId(firebaseUID string) (*models.Users, error)
	AddDefaultCity(defaultCity *models.UserCities) error
	UpdateDeletedUser(user *models.Users, firebaseUID string) error
	GetCitiesByUserId(userID string) ([]*responses.UserCityWithDetails, error)
	IsUserDeleted(userID string) (bool, error)
	GetUserIdByEmail(email string) (string, error)
	SoftDeleteUserGroupForUser(userId string) error
	BaseRepositoryInterface[models.Users]
}

type UserRepository struct {
	BaseRepository[models.Users] // Embed BaseRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: BaseRepository[models.Users]{DB: db}}
}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func (ur *UserRepository) GetUserByEmail(email string) (*models.Users, error) {
	var user models.Users
	userResult := ur.DB.Where("email = ?", email).First(&user)
	if userResult.Error != nil {
		if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user from database: %v", userResult.Error)
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByFirebaseId(firebaseUID string) (*models.Users, error) {
	var user models.Users
	userResult := initializers.Database.Where("Id = ?", firebaseUID).First(&user)
	if userResult.Error != nil {
		if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user from database: %v", userResult.Error)
	}
	return &user, nil
}

func (ur *UserRepository) AddDefaultCity(defaultCity *models.UserCities) error {
	return ur.DB.Create(defaultCity).Error
}

func (ur *UserRepository) UpdateDeletedUser(user *models.Users, firebaseUID string) error {
	user.DeletedAt = nil
	user.Id = firebaseUID //update to new uid
	result := initializers.Database.Save(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to re-reauthenticate user")
	}
	return nil
}

func (ur *UserRepository) GetCitiesByUserId(userID string) ([]*responses.UserCityWithDetails, error) {
	var userCitiesResponse []*responses.UserCityWithDetails
	err := ur.DB.Table("user_cities").
		Select("cities.id AS id, cities.name AS city_name, cities.description AS city_description, cities.image_url AS image_url").
		Joins("INNER JOIN cities ON cities.id = user_cities.city_id").
		Where("user_cities.user_id = ? AND user_cities.deleted_at IS NULL", userID).
		Order("cities.name").
		Find(&userCitiesResponse).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user cities for userID %s: %w", userID, err)
	}
	return userCitiesResponse, nil
}

func (ur *UserRepository) IsUserDeleted(userID string) (bool, error) {
	var count int64
	err := ur.DB.Table("users").
		Where("id = ? AND deleted_at IS NOT NULL", userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if user is deleted for userID %s: %w", userID, err)
	}
	return count > 0, nil
}

func (ur *UserRepository) GetUserIdByEmail(email string) (string, error) {
	var user models.Users

	result := ur.DB.Where("email = ?", email).Find(&user)
	if result.Error != nil {
		return "", fmt.Errorf("UserRepository.GetUserIdsByEmail: Query execution failed, Error: %w", result.Error)
	}
	return user.Id, nil
}

func (ur *UserRepository) SoftDeleteUserGroupForUser(userId string) error {
	// Get the current time
	now := time.Now()

	// Execute the SQL update query to softDelete user_groups for the given user
	err := ur.DB.Model(&models.UserGroup{}).Where("user_id = ? AND deleted_at IS NULL", userId).Update("deleted_at", now).Error

	if err != nil {
		return fmt.Errorf("failed to soft delete user groups for user ID %s: %w", userId, err)
	}

	return nil
}
