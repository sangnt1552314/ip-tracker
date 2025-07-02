package services

import (
	// "log"
	"os"
	"strings"

	"github.com/sangnt1552314/ip-tracker/internal/models"
)

func GetWorldMapText() string {
	worldMap, err := os.ReadFile("assests/map.txt")
    if err != nil {
        return "Error: Could not load world map"
    }
	return string(worldMap)
}

func GetWorldMapDimensions() models.MapDimensions {
	worldMap := GetWorldMapText()

	lines := strings.Split(worldMap, "\n")

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	return models.MapDimensions{
		Width:  maxWidth,
		Height: len(lines),
	}
}

func LatLongToMapPosition(lat, long float64) models.MapPosition {
	dims := GetWorldMapDimensions()

	// Convert latitude (-90 to 90) to Y position (0 to height-1)
	// Note: Y is inverted (0 is top, height-1 is bottom)
	normalizedLat := (90.0 - lat) / 180.0 // 0 to 1
	y := int(normalizedLat * float64(dims.Height))

	// Convert longitude (-180 to 180) to X position (0 to width-1)
	normalizedLong := (long + 180.0) / 360.0 // 0 to 1
	x := int(normalizedLong * float64(dims.Width))

	// Ensure positions are within bounds
	if x < 0 {
		x = 0
	}
	if x >= dims.Width {
		x = dims.Width - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= dims.Height {
		y = dims.Height - 1
	}

	return models.MapPosition{X: x, Y: y}
}