package http

import (
	"wan-system/internal/model"
	"wan-system/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TelegrafController struct {
	Log             *logrus.Logger
	TelegrafUseCase *usecase.TelegrafUseCase
}

func NewTelegrafController(logger *logrus.Logger, useCase *usecase.TelegrafUseCase) *TelegrafController {
	return &TelegrafController{
		Log:             logger,
		TelegrafUseCase: useCase,
	}
}

func (c *TelegrafController) ListAgents(ctx *fiber.Ctx) error {
	responses, err := c.TelegrafUseCase.ListAgents(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.TelegrafAgentResponse]{Data: responses})
}

func (c *TelegrafController) CreateAgent(ctx *fiber.Ctx) error {
	request := new(model.TelegrafAgentRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.TelegrafUseCase.CreateAgent(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TelegrafAgentResponse]{Data: response})
}

func (c *TelegrafController) GenerateConfig(ctx *fiber.Ctx) error {
	config, err := c.TelegrafUseCase.GenerateSnmpConfig(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[string]{Data: config})
}
