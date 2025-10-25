package bidder

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hemantkumar-dev/auction-simulator/internal/model"
)

type BidRequest struct {
	AuctionID    int
	Attributes   model.Attribute
	ResponseCh   chan<- model.Bid
	AuctionStart time.Time
	Ctx          context.Context
}

func NewBidder(id int, aggressiveness float64, latencyMean int) func(BidRequest) {
	return func(req BidRequest) {
		probBid := 0.45 + aggressiveness*0.5
		delay := time.Duration(rand.Intn(latencyMean)+latencyMean/2) * time.Millisecond

		select {
		case <-req.Ctx.Done():
			return
		case <-time.After(delay):
		}

		if req.Ctx.Err() != nil || rand.Float64() > probBid {
			return
		}

		sum := 0.0
		for _, v := range req.Attributes {
			sum += v
		}
		base := sum / float64(len(req.Attributes))
		amount := base*(0.2+aggressiveness) + rand.Float64()*base*0.5

		now := time.Now()
		timeMs := now.Sub(req.AuctionStart).Milliseconds()

		select {
		case <-req.Ctx.Done():
			return
		case req.ResponseCh <- model.Bid{
			BidderID: id,
			Amount:   amount,
			TimeMs:   timeMs,
			Meta:     fmt.Sprintf("latency=%dms", int(delay/time.Millisecond)),
		}:
		case <-time.After(5 * time.Millisecond):
			return
		}
	}
}
