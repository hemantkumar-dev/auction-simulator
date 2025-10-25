package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hemantkumar-dev/auction-simulator/internal/auction"
	"github.com/hemantkumar-dev/auction-simulator/internal/bidder"
	"github.com/hemantkumar-dev/auction-simulator/internal/model"
	"github.com/hemantkumar-dev/auction-simulator/internal/resources"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Prompt user for runtime input
	numBidders := readInt(reader, "Enter number of bidders (e.g. 100): ")
	numAuctions := readInt(reader, "Enter number of concurrent auctions (e.g. 40): ")
	timeoutMs := readInt(reader, "Enter auction timeout in milliseconds (e.g. 800): ")

	vcpu := runtime.NumCPU()
	ramMb := 2048
	outDir := "sample-outputs"

	// ----------------------------
	// Clean up old output files
	os.MkdirAll(outDir, 0755)
	files, err := os.ReadDir(outDir)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
				os.Remove(filepath.Join(outDir, f.Name()))
			}
		}
	}
	// ----------------------------

	rand.Seed(time.Now().UnixNano())
	os.MkdirAll(outDir, 0755)

	fmt.Printf("\nAuctions=%d | Bidders=%d | Timeout=%dms\n", numAuctions, numBidders, timeoutMs)

	// Resource limiter
	limiter := resources.NewLimiter(vcpu, ramMb)

	// Create bidders
	bidders := make([]func(bidder.BidRequest), numBidders)
	for i := 0; i < numBidders; i++ {
		bidders[i] = bidder.NewBidder(i, rand.Float64(), 50+rand.Intn(400))
	}

	results := make(chan model.AuctionResult, numAuctions)
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(numAuctions)

	for a := 0; a < numAuctions; a++ {
		id := a
		go func() {
			defer wg.Done()
			result := auction.RunAuction(id, bidders, timeoutMs, limiter)
			results <- result
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fn := filepath.Join(outDir, fmt.Sprintf("auction_%03d.json", res.AuctionID))
		saveJSON(fn, res)
	}

	end := time.Now()
	total := end.Sub(start)
	fmt.Printf("\nAll auctions complete. Elapsed: %v\n", total)

	saveJSON(filepath.Join(outDir, "summary.json"), map[string]interface{}{
		"start_time": start,
		"end_time":   end,
		"elapsed_ms": total.Milliseconds(),
		"auctions":   numAuctions,
		"bidders":    numBidders,
	})
}

func readInt(reader *bufio.Reader, prompt string) int {
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		value, err := strconv.Atoi(input)
		if err == nil && value > 0 {
			return value
		}
		fmt.Println("Invalid input, please enter a positive integer.")
	}
}

func saveJSON(filename string, v interface{}) {
	f, _ := os.Create(filename)
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
