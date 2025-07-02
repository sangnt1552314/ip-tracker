package app

import (
	// "log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sangnt1552314/ip-tracker/internal/models"
	"github.com/sangnt1552314/ip-tracker/internal/services"
)

type App struct {
	*tview.Application
	ipInfo *models.CloudflareGeoLocationDetail
}

func NewApp() *App {
	app := &App{
		Application: tview.NewApplication(),
		ipInfo:      nil,
	}

	app.EnableMouse(true)

	app.setupBindings()

	// Fetch IP information and store it in the app
	var defaultIpDetail *models.CloudflareGeoLocationDetail
	defaultIpDetail, err := services.GetGeoLocationData()
	if err != nil {
		app.ipInfo = nil
	} else {
		app.ipInfo = defaultIpDetail
	}

	root := tview.NewFlex()
	app.setupLayout(root)

	app.Application.SetRoot(root, true)

	return app
}

func (a *App) Run() error {
	return a.Application.Run()
}

func (a *App) setupBindings() {
	a.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.Stop()
			return nil
		}
		return event
	})
}

func (a *App) setupLayout(root *tview.Flex) {
	root.SetDirection(tview.FlexRow).SetBorder(false)

	menu := a.setupMainMenu()

	content := a.setupMainContent()

	root.AddItem(content, 0, 1, false)
	root.AddItem(menu, 3, 0, false)
}

func (a *App) setupMainMenu() tview.Primitive {
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)
	menuFlex.SetTitle("Options").SetTitleAlign(tview.AlignLeft).SetTitleColor(tcell.ColorWhite)
	menuFlex.SetBorderPadding(0, 0, 0, 0)

	exitButton := tview.NewButton("Exit")
	exitButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))
	exitButton.SetSelectedFunc(func() {
		a.Application.Stop()
	})

	menuFlex.AddItem(exitButton, 9, 0, false)

	return menuFlex
}

func (a *App) setupMainContent() tview.Primitive {
	contentFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	worldmapFlex := tview.NewFlex()
	worldmapFlex.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)
	worldmapFlex.SetTitle("World Map").SetTitleAlign(tview.AlignLeft).SetTitleColor(tcell.ColorWhite)

	// Add the world map to the worldmapFlex
	worldMap := a.createWorldMap()
	worldmapFlex.AddItem(worldMap, 0, 1, false)

	infoFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	infoFlex.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)
	infoFlex.SetTitle("IP Information").SetTitleAlign(tview.AlignLeft).SetTitleColor(tcell.ColorWhite)
	infoFlex.SetBorderPadding(0, 0, 0, 0)

	// Fetch and display IP information
	if a.ipInfo == nil {
		infoFlex.AddItem(tview.NewTextView().SetText("Error fetching IP information"), 0, 1, false)
	} else {
		infoText := tview.NewTextView().SetText(
			"● IP: " + a.ipInfo.ClientIp + "\n" +
				"● Country: " + a.ipInfo.Country + "\n" +
				"● ASN: " + string(a.ipInfo.Asn) + "\n" +
				"● City: " + a.ipInfo.City + "\n" +
				"● Region: " + a.ipInfo.Region + "\n" +
				"● [green]Latitude[white]: " + "[green]" + a.ipInfo.Latitude + "[white]" + "\n" +
				"● [green]Longitude[white]: " + "[green]" + a.ipInfo.Longitude + "[white]" + "\n" +
				"● Colo: " + a.ipInfo.Colo + "\n" +
				"● Postal Code: " + a.ipInfo.PostalCode + "\n" +
				"● HTTP Protocol: " + a.ipInfo.HttpProtocol + "\n" +
				"● ASN Organization: " + a.ipInfo.AsOrganization + "\n" +
				"● Hostname: " + a.ipInfo.Hostname + "\n",
		).
			SetDynamicColors(true)
		infoFlex.AddItem(infoText, 0, 1, false)
	}

	contentFlex.AddItem(worldmapFlex, 0, 8, false)
	contentFlex.AddItem(infoFlex, 0, 2, false)

	return contentFlex
}

// createWorldMap creates a text view with an ASCII world map
func (a *App) createWorldMap() *tview.TextView {
	worldMap := tview.NewTextView()
	worldMap.SetDynamicColors(true)
	worldMap.SetRegions(true)
	worldMap.SetScrollable(true)

	// Detailed ASCII World Map (Flat Projection - Mercator Style)
	mapContent := services.GetWorldMapText()

	if a.ipInfo != nil && a.ipInfo.Latitude != "" && a.ipInfo.Longitude != "" {
		mapContent = a.addLocationMarker(mapContent)
	}

	worldMap.SetText(mapContent)

	return worldMap
}

func (a *App) addLocationMarker(mapContent string) string {
	// Parse latitude and longitude
	lat, err1 := strconv.ParseFloat(a.ipInfo.Latitude, 64)
	long, err2 := strconv.ParseFloat(a.ipInfo.Longitude, 64)

	if err1 != nil || err2 != nil {
		return mapContent // Return original map if parsing fails
	}

	// Convert to map position
	position := services.LatLongToMapPosition(lat, long)

	// Split map into lines
	lines := strings.Split(mapContent, "\n")

	// Ensure we have enough lines
	if position.Y >= len(lines) {
		return mapContent
	}

	// Convert line to rune slice for proper character handling
	runes := []rune(lines[position.Y])

	// Ensure position X is within bounds
	if position.X >= len(runes) {
		return mapContent
	}

	// Replace character at position with marker
	newLine := string(runes[:position.X]) + "[red]●[white]" + string(runes[position.X+1:])
	lines[position.Y] = newLine

	return strings.Join(lines, "\n")
}
