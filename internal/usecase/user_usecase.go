package usecase

import (
	"context"
	"wan-system/internal/entity"
	"wan-system/internal/model"
	"wan-system/internal/repository"
	"wan-system/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserUseCase holds the business logic for user authentication
type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	TokenUtil      *util.TokenUtil
}

func NewUserUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	userRepository *repository.UserRepository,
	tokenUtil *util.TokenUtil,
) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		TokenUtil:      tokenUtil,
	}
}

// Register creates a new user with a hashed password
func (uc *UserUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warnf("Invalid register request: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if user ID already exists
	var count int64
	if err := tx.Model(&entity.User{}).Where("id = ?", request.ID).Count(&count).Error; err != nil {
		uc.Log.Warnf("Failed to count user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if count > 0 {
		uc.Log.Warnf("User %s already exists", request.ID)
		return nil, fiber.ErrConflict
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.Log.Warnf("Failed to hash password: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		ID:       request.ID,
		Password: string(hashed),
		Name:     request.Name,
		Email:    request.Email,
	}

	if err := uc.UserRepository.Create(tx, user); err != nil {
		uc.Log.Warnf("Failed to create user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return toUserResponse(user), nil
}

// Login validates credentials and returns a JWT token
func (uc *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warnf("Invalid login request: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := uc.UserRepository.FindById(tx, user, request.ID); err != nil {
		uc.Log.Warnf("User not found: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		uc.Log.Warnf("Wrong password for user %s", request.ID)
		return nil, fiber.ErrUnauthorized
	}

	jwtToken, err := uc.TokenUtil.GenerateJWT(user)
	if err != nil {
		uc.Log.Warnf("Failed to generate JWT: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	user.Token = jwtToken

	// Persist token to DB
	if err := uc.UserRepository.Update(tx, user); err != nil {
		uc.Log.Warnf("Failed to save token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	resp := toUserResponse(user)
	resp.Token = jwtToken
	return resp, nil
}

// Current returns the currently authenticated user's data
func (uc *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warnf("Invalid current user request: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := uc.UserRepository.FindById(tx, user, request.ID); err != nil {
		uc.Log.Warnf("User not found: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return toUserResponse(user), nil
}

// Logout clears the user's token from the database
func (uc *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warnf("Invalid logout request: %+v", err)
		return false, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := uc.UserRepository.FindById(tx, user, request.ID); err != nil {
		uc.Log.Warnf("User not found on logout: %+v", err)
		return false, fiber.ErrNotFound
	}

	user.Token = ""
	if err := uc.UserRepository.Update(tx, user); err != nil {
		uc.Log.Warnf("Failed to clear token: %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Warnf("Failed to commit logout: %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

// toUserResponse maps a User entity to a UserResponse
func toUserResponse(u *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
