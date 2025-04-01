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

// TrackingHandler handles HTTP requests related to tracking.
type TrackingHandler struct {
	srv *service.TrackingService
	log *zap.SugaredLogger
}

// NewTrackingHandler creates a new TrackingHandler
func NewTrackingHandler(srv *service.TrackingService, log *zap.SugaredLogger) *TrackingHandler {
	return &TrackingHandler{
		srv: srv,
		log: log,
	}
}

func (h *TrackingHandler) TrackEvent(c *fiber.Ctx) error {
	var (
		req *model.TrackingEvent
		err error
	)

	if err = c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err = validation.ValidateStruct(req,
		// Event type
		validation.Field(&req.EventType,
			validation.Required,
			validation.In(
				model.TrackingEventTypeImpression,
				model.TrackingEventTypeClick,
				model.TrackingEventTypeConversion,
			),
		),
		// Line item ID
		validation.Field(&req.LineItemID,
			validation.Required,
		),
		// Timestamp
		validation.Field(&req.Timestamp,
			validation.Required,
		),
		// Placement
		validation.Field(&req.Placement,
			append([]validation.Rule{validation.Required}, internal.PlacementRules...)...,
		),
		// User ID
		validation.Field(&req.UserID,
			validation.Required,
			validation.Length(2, 30),
		),
		// Metadata
		validation.Field(&req.Metadata,
			validation.Required,
		),
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request parameters",
			"details": internal.ValidationErrorJSON(err),
		})
	}

	h.log.Debugw("track event",
		"type", req.EventType, "line_item_id", req.LineItemID, "placement", req.Placement,
		"user_id", req.UserID, "metadata", req.Metadata,
	)

	// Track event
	if err = h.srv.TrackEvent(req); err != nil {

		// Check error is line item not found
		if errors.Is(err, service.ErrLineItemNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Line item not found",
				"details": err.Error(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Unable to track event",
				"details": err.Error(),
			})
		}

	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
	})
}
