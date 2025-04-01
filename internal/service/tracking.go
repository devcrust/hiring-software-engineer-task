package service

import (
	"fmt"
	"sync"

	"sweng-task/internal/model"

	"go.uber.org/zap"
)

// NewTrackingService creates a new TrackingService.
func NewTrackingService(lineItemSrv *LineItemService, log *zap.SugaredLogger) *TrackingService {
	return &TrackingService{
		lineItemSrv: lineItemSrv,
		log:         log,
	}
}

// TrackingService the tracking service.
type TrackingService struct {
	lineItemSrv *LineItemService
	log         *zap.SugaredLogger
	events      []*model.TrackingEvent
	mu          sync.RWMutex
}

// TrackEvent tracks the given event.
func (s *TrackingService) TrackEvent(event *model.TrackingEvent) error {
	var (
		lineItem *model.LineItem
		err      error
	)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieve line item
	if lineItem, err = s.lineItemSrv.GetByID(event.LineItemID); err != nil {
		return fmt.Errorf("unable to retrieve line item: %w", err)
	}

	// Reduce budget based on the bid
	lineItem.Budget = lineItem.Budget - lineItem.Bid

	// Check remaining budget
	if lineItem.Budget <= 0 {
		lineItem.Status = model.LineItemStatusCompleted
	}

	s.log.Infow("line item updated",
		"line_item_id", lineItem.ID, "budget", lineItem.Budget, "status", lineItem.Status,
	)

	// Track event
	s.events = append(s.events, event)

	return nil
}
