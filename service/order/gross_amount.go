package order

import (
	"fmt"
	"strconv"
	"strings"
)

// parseGrossAmount mengubah gross amount (string, misal "10000.00") menjadi int64
func parseGrossAmount(gross string) (int64, error) {
	trimmed := strings.TrimSpace(gross)
	if trimmed == "" {
		return 0, fmt.Errorf("empty gross amount")
	}
	amount, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, err
	}
	return int64(amount), nil
}
