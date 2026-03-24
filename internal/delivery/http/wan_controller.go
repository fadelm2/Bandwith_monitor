package http

import (
	"math"
	"wan-system/internal/model"
	"wan-system/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WanController struct {
	Log     *logrus.Logger
	UseCase *usecase.WanUseCase
}

func NewWanController(logger *logrus.Logger, useCase *usecase.WanUseCase) *WanController {
	return &WanController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *WanController) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": "ok"})
}

func (c *WanController) CreateCapacity(ctx *fiber.Ctx) error {
	request := new(model.WanCapacityRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.CreateCapacity(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WanCapacityResponse]{Data: response})
}

func (c *WanController) UpdateCapacity(ctx *fiber.Ctx) error {
	request := new(model.WanCapacityRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	request.WanID = ctx.Params("wanId")
	response, err := c.UseCase.UpdateCapacity(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WanCapacityResponse]{Data: response})
}

func (c *WanController) DeleteCapacity(ctx *fiber.Ctx) error {
	wanID := ctx.Params("wanId")
	if err := c.UseCase.DeleteCapacity(ctx.UserContext(), wanID); err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

func (c *WanController) GetCapacity(ctx *fiber.Ctx) error {
	wanID := ctx.Params("wanId")
	response, err := c.UseCase.GetCapacity(ctx.UserContext(), wanID)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WanCapacityResponse]{Data: response})
}

func (c *WanController) ListCapacity(ctx *fiber.Ctx) error {
	responses, err := c.UseCase.ListCapacity(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]*model.WanCapacityResponse]{Data: responses})
}

func (c *WanController) BulkUpdateCapacity(ctx *fiber.Ctx) error {
	var request []model.WanCapacityRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.BulkUpdateCapacity(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]*model.WanCapacityResponse]{Data: response})
}

func (c *WanController) SearchTraffic(ctx *fiber.Ctx) error {
	request := &model.SearchTrafficRequest{
		WanID: ctx.Query("wan_id", ""),
		Page:  ctx.QueryInt("page", 1),
		Size:  ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.SearchTraffic(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]*model.WanTrafficResponse]{
		Data:   responses,
		Paging: paging,
	})
}
