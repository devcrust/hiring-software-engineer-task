package handler

import (
	"errors"

	"sweng-task/internal"
	"sweng-task/internal/model"

	"sweng-task/internal/service"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	// LineItemDetailsRoute the line item details route name to use.
	LineItemDetailsRoute = "line_item.details"
)

// LineItemHandler handles HTTP requests related to line items
type LineItemHandler struct {
	service *service.LineItemService
	log     *zap.SugaredLogger
}

// NewLineItemHandler creates a new LineItemHandler
func NewLineItemHandler(service *service.LineItemService, log *zap.SugaredLogger) *LineItemHandler {
	return &LineItemHandler{
		service: service,
		log:     log,
	}
}

// Create handles the creation of a new line item
func (h *LineItemHandler) Create(c *fiber.Ctx) error {
	var (
		input model.LineItemCreate
		err   error
	)
	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err = validation.ValidateStruct(&input,
		// Name
		validation.Field(&input.Name,
			append([]validation.Rule{validation.Required}, internal.NameRules...)...,
		),
		// Advertiser ID
		validation.Field(&input.AdvertiserID,
			append([]validation.Rule{validation.Required}, internal.AdvertiserIDRules...)...,
		),
		// Bid
		validation.Field(&input.Bid,
			validation.Required,
			validation.Min(0.0),
			validation.Max(50.0),
		),
		// Budget
		validation.Field(&input.Budget,
			validation.Required,
			validation.Min(input.Bid),
			validation.Max(10_000.0),
		),
		// Placement
		validation.Field(&input.Placement,
			append([]validation.Rule{validation.Required}, internal.PlacementRules...)...,
		),
		// Categories
		validation.Field(&input.Categories,
			validation.Required,
			validation.Each(internal.CategoryRules...),
		),
		// Keywords
		validation.Field(&input.Keywords,
			validation.Required,
			validation.Each(internal.KeywordRules...),
		),
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request parameters",
			"details": internal.ValidationErrorJSON(err),
		})
	}

	lineItem, err := h.service.Create(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to create line item",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(lineItem)
}

// GetByID handles retrieving a line item by ID
func (h *LineItemHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Missing line item ID",
		})
	}

	lineItem, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrLineItemNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    fiber.StatusNotFound,
				"message": "Line item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to retrieve line item",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(lineItem)
}

// GetAll handles retrieving all line items with optional filtering
func (h *LineItemHandler) GetAll(c *fiber.Ctx) error {
	var (
		req = struct {
			AdvertiserID string
			Placement    string
		}{
			AdvertiserID: c.Params("advertiser_id"),
			Placement:    c.Query("placement"),
		}
		err error
	)

	// Validate request
	if err = validation.ValidateStruct(&req,
		// Advertiser ID
		validation.Field(&req.AdvertiserID,
			append([]validation.Rule{validation.Required}, internal.AdvertiserIDRules...)...,
		),
		// Placement
		validation.Field(&req.Placement,
			append([]validation.Rule{validation.Required}, internal.PlacementRules...)...,
		),
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request parameters",
			"details": internal.ValidationErrorJSON(err),
		})
	}

	lineItems, err := h.service.GetAll(req.AdvertiserID, req.Placement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to retrieve line items",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(lineItems)
}
