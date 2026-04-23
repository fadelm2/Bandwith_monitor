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

// Health check
// @Summary Health check
// @Description Check system status
// @Tags System
// @Produce json
// @Success 200 {object} map[string]string
// @Router /internal/health [get]
func (c *WanController) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": "ok"})
}

// CreateCapacity
// @Summary Create WAN capacity
// @Tags Capacity
// @Accept json
// @Produce json
// @Param request body model.WanCapacityRequest true "Capacity Info"
// @Success 200 {object} model.WebResponse[model.WanCapacityResponse]
// @Security bearerAuth
// @Router /internal/capacity [post]
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

// UpdateCapacity
// @Summary Update WAN capacity
// @Tags Capacity
// @Accept json
// @Produce json
// @Param wanId path string true "WAN ID"
// @Param request body model.WanCapacityRequest true "Capacity Info"
// @Success 200 {object} model.WebResponse[model.WanCapacityResponse]
// @Security bearerAuth
// @Router /internal/capacity/{wanId} [put]
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

// DeleteCapacity
// @Summary Delete WAN capacity
// @Tags Capacity
// @Param wanId path string true "WAN ID"
// @Success 200 {object} model.WebResponse[bool]
// @Security bearerAuth
// @Router /internal/capacity/{wanId} [delete]
func (c *WanController) DeleteCapacity(ctx *fiber.Ctx) error {
	wanID := ctx.Params("wanId")
	if err := c.UseCase.DeleteCapacity(ctx.UserContext(), wanID); err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

// GetCapacity
// @Summary Get WAN capacity
// @Tags Capacity
// @Produce json
// @Param wanId path string true "WAN ID"
// @Success 200 {object} model.WebResponse[model.WanCapacityResponse]
// @Security bearerAuth
// @Router /internal/capacity/{wanId} [get]
func (c *WanController) GetCapacity(ctx *fiber.Ctx) error {
	wanID := ctx.Params("wanId")
	response, err := c.UseCase.GetCapacity(ctx.UserContext(), wanID)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.WanCapacityResponse]{Data: response})
}

// ListCapacity
// @Summary List all WAN capacities
// @Tags Capacity
// @Produce json
// @Success 200 {object} model.WebResponse[[]model.WanCapacityResponse]
// @Security bearerAuth
// @Router /internal/capacity [get]
func (c *WanController) ListCapacity(ctx *fiber.Ctx) error {
	responses, err := c.UseCase.ListCapacity(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]*model.WanCapacityResponse]{Data: responses})
}

// BulkUpdateCapacity
// @Summary Bulk update WAN capacities
// @Tags Capacity
// @Accept json
// @Produce json
// @Param request body []model.WanCapacityRequest true "List of Capacity Info"
// @Success 200 {object} model.WebResponse[[]model.WanCapacityResponse]
// @Security bearerAuth
// @Router /internal/capacity/bulk [put]
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

// SearchTraffic
// @Summary Search traffic logs
// @Tags Traffic
// @Produce json
// @Param wan_id query string false "WAN ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} model.WebResponse[[]model.WanTrafficResponse]
// @Security bearerAuth
// @Router /internal/traffic [get]
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

// GetAlerts
// @Summary Get high utilization alerts (aggregated)
// @Tags Traffic
// @Produce json
// @Success 200 {object} model.WebResponse[[]model.WanAlertResponse]
// @Security bearerAuth
// @Router /internal/alerts [get]
func (c *WanController) GetAlerts(ctx *fiber.Ctx) error {
	responses, err := c.UseCase.GetAlerts(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.WanAlertResponse]{Data: responses})
}
