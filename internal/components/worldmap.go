package components

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sangnt1552314/ip-tracker/internal/services"
	"github.com/sangnt1552314/ip-tracker/pkg/utils"
)

// MapPoint represents a point on the map with location and color
type MapPoint struct {
	Latitude  float64
	Longitude float64
	Color     string // Changed from tcell.Color to string
	Symbol    rune
}

// WorldMapWidget represents a custom draggable world map component
type WorldMapWidget struct {
	*tview.Box
	mapContent string
	offsetX    int
	offsetY    int
	isDragging bool
	lastMouseX int
	lastMouseY int
	points     []MapPoint
}

// NewWorldMapWidget creates a new draggable world map widget
func NewWorldMapWidget() *WorldMapWidget {
	widget := &WorldMapWidget{
		Box:        tview.NewBox(),
		mapContent: services.GetWorldMapText(),
		offsetX:    0,
		offsetY:    0,
		isDragging: false,
		points:     make([]MapPoint, 0),
	}

	widget.SetInputCapture(widget.handleInput)
	widget.SetMouseCapture(widget.handleMouse)

	return widget
}

// Draw renders the world map with current offset
func (w *WorldMapWidget) Draw(screen tcell.Screen) {
	w.Box.DrawForSubclass(screen, w)
	x, y, width, height := w.GetInnerRect()

	// Get the base map content with points
	baseMap := services.GetWorldMapText()
	mapWithPoints := w.addPointsToMap(baseMap)
	lines := strings.Split(mapWithPoints, "\n")

	for i := 0; i < height; i++ {
		if i+w.offsetY < 0 || i+w.offsetY >= len(lines) {
			continue
		}

		line := lines[i+w.offsetY]
		runes := []rune(line)

		currentColor := tcell.ColorWhite
		screenX := 0

		for j := 0; j < len(runes) && screenX < width; {
			if j+w.offsetX < 0 {
				j++
				continue
			}

			// Check bounds before accessing runes array
			if j+w.offsetX >= len(runes) {
				break
			}

			char := runes[j+w.offsetX]

			// Check for color tag
			if char == '[' {
				// Find the end of color tag
				tagEnd := -1
				for k := j + w.offsetX + 1; k < len(runes); k++ {
					if runes[k] == ']' {
						tagEnd = k
						break
					}
				}

				if tagEnd != -1 {
					colorTag := string(runes[j+w.offsetX : tagEnd+1])
					// Extract color name from tag and convert to tcell color
					colorName := colorTag[1 : len(colorTag)-1] // Remove [ and ]
					currentColor = w.getColorByName(colorName)

					// Skip the entire color tag
					j += (tagEnd - (j + w.offsetX)) + 1
					continue
				}
			}

			// Render the character with current color
			if screenX >= 0 && screenX < width {
				style := tcell.StyleDefault.Foreground(currentColor)
				screen.SetContent(x+screenX, y+i, char, nil, style)
				screenX++
			}

			j++
		}
	}
}

// handleInput handles keyboard input for the world map
func (w *WorldMapWidget) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		w.offsetY = utils.Max(w.offsetY-1, 0)
		return nil
	case tcell.KeyDown:
		w.offsetY++
		return nil
	case tcell.KeyLeft:
		w.offsetX = utils.Max(w.offsetX-1, 0)
		return nil
	case tcell.KeyRight:
		w.offsetX++
		return nil
	case tcell.KeyHome:
		w.offsetX = 0
		w.offsetY = 0
		return nil
	}

	return event
}

// handleMouse handles mouse events for dragging
func (w *WorldMapWidget) handleMouse(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	x, y := event.Position()

	switch action {
	case tview.MouseLeftDown:
		w.isDragging = true
		w.lastMouseX = x
		w.lastMouseY = y
		return tview.MouseConsumed, nil

	case tview.MouseMove:
		if w.isDragging {
			deltaX := w.lastMouseX - x
			deltaY := w.lastMouseY - y

			w.offsetX = utils.Max(w.offsetX+deltaX, 0)
			w.offsetY = utils.Max(w.offsetY+deltaY, 0)

			w.lastMouseX = x
			w.lastMouseY = y
		}
		return tview.MouseConsumed, nil

	case tview.MouseLeftUp:
		w.isDragging = false
		return tview.MouseConsumed, nil

	case tview.MouseScrollUp:
		w.offsetY = utils.Max(w.offsetY-3, 0)
		return tview.MouseConsumed, nil

	case tview.MouseScrollDown:
		w.offsetY += 3
		return tview.MouseConsumed, nil
	}

	return action, event
}

// AddPoint adds a point to the map
func (w *WorldMapWidget) AddPoint(lat, long float64, color string, symbol rune) {
	point := MapPoint{
		Latitude:  lat,
		Longitude: long,
		Color:     color,
		Symbol:    symbol,
	}
	w.points = append(w.points, point)
}

// SetPoints sets all points on the map
func (w *WorldMapWidget) SetPoints(points []MapPoint) {
	w.points = points
}

// ClearPoints removes all points from the map
func (w *WorldMapWidget) ClearPoints() {
	w.points = make([]MapPoint, 0)
}

// addPointsToMap adds all points to the map content
func (w *WorldMapWidget) addPointsToMap(mapContent string) string {
	lines := strings.Split(mapContent, "\n")

	for _, point := range w.points {
		position := services.LatLongToMapPosition(point.Latitude, point.Longitude)

		if position.Y >= len(lines) || position.Y < 0 {
			continue
		}

		runes := []rune(lines[position.Y])
		if position.X >= len(runes) || position.X < 0 {
			continue
		}

		// Convert string color to color tag
		colorTag := "[" + point.Color + "]"

		newLine := string(runes[:position.X]) + colorTag + string(point.Symbol) + "[white]" + string(runes[position.X+1:])
		lines[position.Y] = newLine
	}

	return strings.Join(lines, "\n")
}

// getColorByName converts color name string to tcell.Color
func (w *WorldMapWidget) getColorByName(colorName string) tcell.Color {
	// Use tcell's built-in color name mapping
	color, ok := tcell.ColorNames[colorName]
	if ok {
		return color
	}
	// Default to white if color name not found
	return tcell.ColorWhite
}
