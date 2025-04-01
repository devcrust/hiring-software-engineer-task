package service

import (
	"errors"
	"sync"
	"time"

	"sweng-task/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Errors
var (
	ErrLineItemNotFound = errors.New("line item not found")
)

// LineItemService provides operations for line items
type LineItemService struct {
	items map[string]*model.LineItem
	mu    sync.RWMutex
	log   *zap.SugaredLogger
}

// NewLineItemService creates a new LineItemService
func NewLineItemService(log *zap.SugaredLogger) *LineItemService {
	return &LineItemService{
		items: make(map[string]*model.LineItem),
		log:   log,
	}
}

// Create creates a new line item
func (s *LineItemService) Create(item model.LineItemCreate) (*model.LineItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	lineItem := &model.LineItem{
		ID:           "li_" + uuid.New().String(),
		Name:         item.Name,
		AdvertiserID: item.AdvertiserID,
		Bid:          item.Bid,
		Budget:       item.Budget,
		Placement:    item.Placement,
		Categories:   item.Categories,
		Keywords:     item.Keywords,
		Status:       model.LineItemStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.items[lineItem.ID] = lineItem
	s.log.Infow("Line item created",
		"id", lineItem.ID,
		"name", lineItem.Name,
		"advertiser_id", lineItem.AdvertiserID,
		"placement", lineItem.Placement,
	)

	return lineItem, nil
}

// GetByID retrieves a line item by ID
func (s *LineItemService) GetByID(id string) (*model.LineItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, ErrLineItemNotFound
	}

	return item, nil
}

// GetAll retrieves all line items, optionally filtered by advertiser ID and placement
func (s *LineItemService) GetAll(advertiserID, placement string) ([]*model.LineItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.LineItem

	for _, item := range s.items {
		if advertiserID != "" && item.AdvertiserID != advertiserID {
			continue
		}

		if placement != "" && item.Placement != placement {
			continue
		}

		result = append(result, item)
	}

	return result, nil
}

// FindMatchingLineItems finds line items matching the given placement and filters
// This method will be used by the AdService when implementing the ad selection logic
func (s *LineItemService) FindMatchingLineItems(placement, category, keyword string) ([]*model.LineItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.LineItem

	for _, item := range s.items {
		// Skip items not matching the placement or not active
		if item.Placement != placement || item.Status != model.LineItemStatusActive {
			continue
		}

		// Apply category filter if specified
		if category != "" {
			categoryFound := false
			for _, cat := range item.Categories {
				if cat == category {
					categoryFound = true
					break
				}
			}
			if !categoryFound {
				continue
			}
		}

		// Apply keyword filter if specified
		if keyword != "" {
			keywordFound := false
			for _, kw := range item.Keywords {
				if kw == keyword {
					keywordFound = true
					break
				}
			}
			if !keywordFound {
				continue
			}
		}

		result = append(result, item)
	}

	return result, nil
}
