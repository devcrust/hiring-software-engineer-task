package model

import (
	"time"
)

// LineItemStatus represents the status of a line item
type LineItemStatus string

const (
	LineItemStatusActive    LineItemStatus = "active"
	LineItemStatusPaused    LineItemStatus = "paused"
	LineItemStatusCompleted LineItemStatus = "completed"
)

// LineItem represents an advertisement with associated bid information
type LineItem struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	AdvertiserID string         `json:"advertiser_id"`
	Bid          float64        `json:"bid"`
	Budget       float64        `json:"budget"`
	Placement    string         `json:"placement"`
	Categories   []string       `json:"categories,omitempty"`
	Keywords     []string       `json:"keywords,omitempty"`
	Status       LineItemStatus `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// LineItemCreate represents the data needed to create a new line item
type LineItemCreate struct {
	Name         string   `json:"name"`
	AdvertiserID string   `json:"advertiser_id"`
	Bid          float64  `json:"bid"`
	Budget       float64  `json:"budget"`
	Placement    string   `json:"placement"`
	Categories   []string `json:"categories,omitempty"`
	Keywords     []string `json:"keywords,omitempty"`
}
