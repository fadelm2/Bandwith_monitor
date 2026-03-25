package http

import (
	"wan-system/internal/delivery/http/middleware"
	"wan-system/internal/model"
	"wan-system/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// UserController handles HTTP requests related to user authentication
type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(log *logrus.Logger, useCase *usecase.UserUseCase) *UserController {
	return &UserController{Log: log, UseCase: useCase}
}

// Register handles POST /api/auth/register
// @Summary Register a new user
// @Description Register a new user with ID, password, name, and email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RegisterUserRequest true "Registration Info"
// @Success 201 {object} model.WebResponse[model.UserResponse]
// @Failure 400 {object} model.WebResponse[string]
// @Router /api/auth/register [post]
func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse register body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Login handles POST /api/auth/login
// @Summary Login and get JWT
// @Description Login with ID and password to receive a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginUserRequest true "Login Credentials"
// @Success 200 {object} model.WebResponse[model.UserResponse]
// @Failure 401 {object} model.WebResponse[string]
// @Router /api/auth/login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse login body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user: %+v", err)
		return err
	}

	// Also set an HTTP-only cookie for browser clients
	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    response.Token,
		Path:     "/",
		MaxAge:   5 * 3600, // 5 hours
		HTTPOnly: true,
		Secure:   false,
	})

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Current handles GET /api/auth/current (protected)
// @Summary Get current user
// @Description Get currently authenticated user information
// @Tags Auth
// @Produce json
// @Success 200 {object} model.WebResponse[model.UserResponse]
// @Failure 401 {object} model.WebResponse[string]
// @Security bearerAuth
// @Router /api/auth/current [get]
func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{ID: auth.ID}
	response, err := c.UseCase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get current user: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Logout handles POST /api/auth/logout (protected)
// @Summary Logout user
// @Description Logout current user and clear token cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} model.WebResponse[bool]
// @Security bearerAuth
// @Router /api/auth/logout [post]
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.LogoutUserRequest{ID: auth.ID}
	response, err := c.UseCase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to logout user: %+v", err)
		return err
	}

	// Clear cookie
	ctx.Cookie(&fiber.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	})

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}
