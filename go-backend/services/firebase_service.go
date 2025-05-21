package services

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/initializers"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/repositories"
)

type FirebaseServiceInterface interface {
	GetUserByFirebaseId(ctx *gin.Context) (*models.Users, error)
	GetAndVerifyIDToken(ctx *gin.Context) (string, error)
	GetUserByUID(firebaseUID string) (*auth.UserRecord, error)
}

type FirebaseService struct {
	UserRepo repositories.UserRepositoryInterface
	client   *auth.Client
}

func NewFirebaseService(repo repositories.UserRepositoryInterface) (*FirebaseService, error) {
	client, err := initializers.FirebaseApp.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase Auth client")
	}
	return &FirebaseService{client: client, UserRepo: repo}, err
}

var _ FirebaseServiceInterface = (*FirebaseService)(nil)

func (f *FirebaseService) GetUserByFirebaseId(ginCtx *gin.Context) (*models.Users, error) {
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user (UserService.GetUserByFirebaseId) : no firebaseUID in context")
	}
	user, err := f.UserRepo.GetUserByFirebaseId(firebaseUID.(string))
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUserByFirebaseId: %v", err)
	}
	return user, nil
}

func (f *FirebaseService) GetAndVerifyIDToken(ginCtx *gin.Context) (string, error) {
	jwt := ginCtx.Request.Header.Get("Authorization")
	if jwt == "" {
		return "", fmt.Errorf("no idToken provided")
	}
	if len(jwt) > 7 && jwt[:7] == "Bearer " {
		jwt = jwt[7:]
	}
	token, err := f.client.VerifyIDToken(context.Background(), jwt)
	if err != nil {
		return "", fmt.Errorf("failed to verify jwt token")
	}
	return token.UID, nil
}

func (f *FirebaseService) GetUserByUID(firebaseUID string) (*auth.UserRecord, error) {
	userRecord, err := f.client.GetUser(context.Background(), firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firebase user details: %v", err)
	}
	return userRecord, nil
}
