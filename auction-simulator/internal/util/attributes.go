package util

import (
	"fmt"
	"math/rand"

	"github.com/hemantkumar-dev/auction-simulator/internal/model"
)

// GenerateAttributes creates n random float64 attributes.
func GenerateAttributes(n int) model.Attribute {
	attr := make(model.Attribute)
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("attr_%02d", i+1)
		attr[key] = rand.Float64() * 100
	}
	return attr
}
