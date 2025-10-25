package model

import "time"

// Attribute represents auction item attributes
type Attribute map[string]float64

// Bid represents a bid in the auction
type Bid struct {
	BidderID int
	Amount   float64
	TimeMs   int64
	Meta     string
}

// AuctionResult represents the outcome of an auction
type AuctionResult struct {
	AuctionID    int
	Attributes   Attribute
	Bids         []Bid
	Winner       *Bid
	StartTime    time.Time
	EndTime      time.Time
	DurationMs   int64
	TimeoutMs    int64
	TotalBidders int
}
