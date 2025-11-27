package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetTime() time.Time {
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		return time.Now().UTC().Add(-3 * time.Hour)
	}
	return time.Now().In(loc)
}

// ParseTimeFrame parses a timeframe string (e.g., "24h", "7d", "30d") and returns the timestamp representing that duration ago
func ParseTimeFrame(timeFrame string) (time.Time, error) {
	if timeFrame == "" {
		timeFrame = "24h"
	}

	timeFrame = strings.ToLower(strings.TrimSpace(timeFrame))

	// Extract numeric part and unit
	var value int
	var unit string

	if len(timeFrame) < 2 {
		return time.Time{}, fmt.Errorf("invalid timeframe format: %s", timeFrame)
	}

	// Parse based on last character (unit)
	lastChar := timeFrame[len(timeFrame)-1:]
	numericPart := timeFrame[:len(timeFrame)-1]

	var err error
	value, err = strconv.Atoi(numericPart)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timeframe value: %s", timeFrame)
	}

	if value <= 0 {
		return time.Time{}, fmt.Errorf("timeframe value must be positive: %s", timeFrame)
	}

	unit = lastChar

	// Calculate duration based on unit
	var duration time.Duration
	switch unit {
	case "h": // hours
		duration = time.Duration(value) * time.Hour
	case "d": // days
		duration = time.Duration(value) * 24 * time.Hour
	case "w": // weeks
		duration = time.Duration(value) * 7 * 24 * time.Hour
	case "m": // months (approximated as 30 days)
		duration = time.Duration(value) * 30 * 24 * time.Hour
	default:
		return time.Time{}, fmt.Errorf("invalid timeframe unit: %s (valid: h, d, w, m)", unit)
	}

	return time.Now().Add(-duration), nil
}
