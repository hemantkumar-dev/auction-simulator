package auction

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/hemantkumar-dev/auction-simulator/internal/bidder"
	"github.com/hemantkumar-dev/auction-simulator/internal/model"
	"github.com/hemantkumar-dev/auction-simulator/internal/resources"
	"github.com/hemantkumar-dev/auction-simulator/internal/util"
)

func RunAuction(
	auctionID int,
	bidders []func(bidder.BidRequest),
	timeoutMs int,
	limiter *resources.Limiter,
) model.AuctionResult {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	start := time.Now()
	attributes := util.GenerateAttributes(20)

	bidCh := make(chan model.Bid, len(bidders))
	var wg sync.WaitGroup
	wg.Add(len(bidders))

	// Launch bidders safely
	for _, bFunc := range bidders {
		go func(bf func(bidder.BidRequest)) {
			defer wg.Done()
			limiter.Acquire()
			defer limiter.Release()

			req := bidder.BidRequest{
				AuctionID:    auctionID,
				Attributes:   attributes,
				ResponseCh:   bidCh,
				AuctionStart: start,
				Ctx:          ctx,
			}

			// Recover panic in any bidder
			defer func() {
				if r := recover(); r != nil {
					// ignore failed bidder
				}
			}()

			bf(req)
		}(bFunc)
	}

	// Close channel when all bidders finish or context times out
	go func() {
		wg.Wait()
		select {
		case <-ctx.Done():
			// already timed out, no need to close (reading side will exit)
		default:
			close(bidCh)
		}
	}()

	var bids []model.Bid
	collectDone := false

	for !collectDone {
		select {
		case b, ok := <-bidCh:
			if !ok {
				collectDone = true
			} else {
				bids = append(bids, b)
			}
		case <-ctx.Done():
			collectDone = true
		}
	}

	end := time.Now()

	var winner *model.Bid
	if len(bids) > 0 {
		sort.SliceStable(bids, func(i, j int) bool {
			if bids[i].Amount == bids[j].Amount {
				return bids[i].TimeMs < bids[j].TimeMs
			}
			return bids[i].Amount > bids[j].Amount
		})
		winner = &bids[0]
	}

	return model.AuctionResult{
		AuctionID:    auctionID,
		Attributes:   attributes,
		Bids:         bids,
		Winner:       winner,
		StartTime:    start,
		EndTime:      end,
		DurationMs:   end.Sub(start).Milliseconds(),
		TimeoutMs:    timeoutMs,
		TotalBidders: len(bidders),
	}
}
