package utils

import (
	"fmt"
	"strings"
)

func FormatRupiah(amount float64) string {
	intAmount := int64(amount)

	isNegative := false
	if intAmount < 0 {
		isNegative = true
		intAmount = -intAmount
	}

	str := fmt.Sprintf("%d", intAmount)
	n := len(str)

	if n <= 3 {
		if isNegative {
			return fmt.Sprintf("-Rp %s", str)
		}
		return fmt.Sprintf("Rp %s", str)
	}

	var builder strings.Builder
	remainder := n % 3

	if remainder > 0 {
		builder.WriteString(str[:remainder])
		builder.WriteString(".")
	}

	for i := remainder; i < n; i += 3 {
		builder.WriteString(str[i : i+3])
		if i+3 < n {
			builder.WriteString(".")
		}
	}

	result := builder.String()
	if isNegative {
		return fmt.Sprintf("-Rp %s", result)
	}

	return fmt.Sprintf("Rp %s", result)
}
