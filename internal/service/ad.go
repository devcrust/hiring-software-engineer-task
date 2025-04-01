package service

import (
	"fmt"
	"sort"

	"sweng-task/internal/model"

	"go.uber.org/zap"
)

const (
	// CategoryWeight weights the category score compared to the actual bid.
	// A higher weight means less importance for the actual bid.
	CategoryWeight = 0.5
)

// NewAdService creates a new AdService
func NewAdService(lineItemSrv *LineItemService, log *zap.SugaredLogger) *AdService {
	return &AdService{
		lineItemSrv: lineItemSrv,
		log:         log,
		categoryScore: map[string]float64{
			"electronics": 1.5,
			"sale":        1.2,
			"fashion":     1.05,
			"travel":      0.9,
			"finance":     0.55,
		},
	}
}

// AdService the ad service.
type AdService struct {
	lineItemSrv   *LineItemService
	log           *zap.SugaredLogger
	categoryScore map[string]float64
}

// GetWinningAds returns the winning ads based on the given parameters.
func (s *AdService) GetWinningAds(placement, category, keyword string) ([]*model.Ad, error) {
	var (
		result    []*model.Ad
		lineItems []*model.LineItem
		lineItem  *model.LineItem
		score     float64
		err       error
	)

	// Retrieve matching line items
	if lineItems, err = s.lineItemSrv.FindMatchingLineItems(placement, category, keyword); err != nil {
		return nil, fmt.Errorf("unable to retrieve winning ads: %w", err)
	}

	// Check line items available
	if len(lineItems) == 0 {
		return nil, nil
	}

	// Initialise result with a suitable capacity
	result = make([]*model.Ad, 0, len(lineItems))

	// List line items
	for _, lineItem = range lineItems {

		// Reset score for current line item
		score = 1

		// Filter line item if budget is exceeded
		if lineItem.Budget < lineItem.Bid {
			continue
		}

		// Calculate score based on categories
		for _, itemCategory := range lineItem.Categories {
			if v, exists := s.categoryScore[itemCategory]; exists {
				score += v
			}
		}

		// Calculate score based on category weight
		score = lineItem.Bid * (1 + CategoryWeight*score)

		// Add line item to result
		result = append(result, &model.Ad{
			ID:           lineItem.ID,
			Name:         lineItem.Name,
			AdvertiserID: lineItem.AdvertiserID,
			Bid:          lineItem.Bid,
			Placement:    lineItem.Placement,
			ServeURL:     "", // assembled in the route handler, see: AdHandler.GetWinningAds
			Score:        score,
		})
	}

	// Sort line items based on the score
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result, nil
}
