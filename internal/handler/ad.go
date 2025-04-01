package handler

import (
	"sweng-task/internal"
	"sweng-task/internal/model"
	"sweng-task/internal/service"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AdHandler handles HTTP requests related to ads.
type AdHandler struct {
	srv *service.AdService
	log *zap.SugaredLogger
}

// NewAdHandler creates a new AdHandler
func NewAdHandler(srv *service.AdService, log *zap.SugaredLogger) *AdHandler {
	return &AdHandler{
		srv: srv,
		log: log,
	}
}

func (h *AdHandler) GetWinningAds(c *fiber.Ctx) error {
	var (
		req = struct {
			Placement string
			Category  string
			Keyword   string
		}{
			Placement: c.Query("placement"),
			Category:  c.Query("category"),
			Keyword:   c.Query("keyword"),
		}
		result      []*model.Ad
		item        *model.Ad
		lineItemUrl string
		err         error
	)

	// Validate request
	if err = validation.ValidateStruct(&req,
		// Placement
		validation.Field(&req.Placement,
			append([]validation.Rule{validation.Required}, internal.PlacementRules...)...,
		),
		// Category
		validation.Field(&req.Category,
			append([]validation.Rule{validation.Required}, internal.CategoryRules...)...,
		),
		// Keyword
		validation.Field(&req.Keyword,
			append([]validation.Rule{validation.Required}, internal.PlacementRules...)...,
		),
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request parameters",
			"details": internal.ValidationErrorJSON(err),
		})
	}

	h.log.Debugw("fetching winning ads",
		"placement", req.Placement, "category", req.Category, "keyword", req.Keyword,
	)

	// Retrieve winning ads
	if result, err = h.srv.GetWinningAds(req.Placement, req.Category, req.Keyword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to retrieve ads",
			"details": err.Error(),
		})
	}

	// Check winning ads
	if len(result) == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
	}

	// List ads
	for _, item = range result {
		// Assemble line item details route
		if lineItemUrl, err = c.GetRouteURL(LineItemDetailsRoute, fiber.Map{
			"id": item.ID,
		}); err == nil {
			item.ServeURL = lineItemUrl
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Unable to assemble line item details URL",
				"details": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
