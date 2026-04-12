package conversion

import (
	"fmt"
	"strconv"
	"strings"
)

func ParsePriceToInt64(val string) (int64, error) {
	val = strings.TrimSpace(val)
	if val == "" {
		return 0, fmt.Errorf("empty price value")
	}

	// Split integer and decimal parts
	parts := strings.Split(val, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid price format: %s", val)
	}

	intPart := parts[0]
	decPart := "00"

	if len(parts) == 2 {
		decPart = parts[1]
	}

	if len(decPart) == 1 {
		decPart += "0"
	} else if len(decPart) > 2 {
		decPart = decPart[:2]
		// todo: create better logic for decpart
	}

	combined := intPart + decPart

	result, err := strconv.ParseInt(combined, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return result, nil
}

func FormatInt64ToPrice(val int64) string {
	sign := ""
	if val < 0 {
		sign = "-"
		val = -val
	}

	intPart := val / 100
	decPart := val % 100

	return fmt.Sprintf("%s%d.%02d", sign, intPart, decPart)
}
