package app

import (
	// "log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sangnt1552314/ip-tracker/internal/components"
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
			"[yellow]●[white] [cyan]IP[white]: [lightblue]" + a.ipInfo.ClientIp + "[white]\n" +
				"[yellow]●[white] [cyan]Country[white]: [lightgreen]" + a.ipInfo.Country + "[white]\n" +
				"[yellow]●[white] [cyan]ASN[white]: [orange]" + strconv.Itoa(a.ipInfo.Asn) + "[white]\n" +
				"[yellow]●[white] [cyan]City[white]: [lightgreen]" + a.ipInfo.City + "[white]\n" +
				"[yellow]●[white] [cyan]Region[white]: [lightgreen]" + a.ipInfo.Region + "[white]\n" +
				"[yellow]●[white] [cyan]Latitude[white]: [green]" + a.ipInfo.Latitude + "[white]\n" +
				"[yellow]●[white] [cyan]Longitude[white]: [green]" + a.ipInfo.Longitude + "[white]\n" +
				"[yellow]●[white] [cyan]Colo[white]: [magenta]" + a.ipInfo.Colo + "[white]\n" +
				"[yellow]●[white] [cyan]Postal Code[white]: [orange]" + a.ipInfo.PostalCode + "[white]\n" +
				"[yellow]●[white] [cyan]HTTP Protocol[white]: [lightblue]" + a.ipInfo.HttpProtocol + "[white]\n" +
				"[yellow]●[white] [cyan]ASN Organization[white]: [lightcyan]" + a.ipInfo.AsOrganization + "[white]\n" +
				"[yellow]●[white] [cyan]Hostname[white]: [lightblue]" + a.ipInfo.Hostname + "[white]\n",
		).SetDynamicColors(true)
		infoFlex.AddItem(infoText, 0, 1, false)
	}

	contentFlex.AddItem(worldmapFlex, 0, 8, false)
	contentFlex.AddItem(infoFlex, 0, 2, false)

	return contentFlex
}

// createWorldMap creates a draggable world map widget
func (a *App) createWorldMap() tview.Primitive {
	worldMap := components.NewWorldMapWidget()

	// Add IP location marker if IP info is available
	if a.ipInfo != nil && a.ipInfo.Latitude != "" && a.ipInfo.Longitude != "" {
		lat, err1 := strconv.ParseFloat(a.ipInfo.Latitude, 64)
		long, err2 := strconv.ParseFloat(a.ipInfo.Longitude, 64)

		if err1 == nil && err2 == nil {
			worldMap.AddPoint(lat, long, "red", '●')
		}
	}

	return worldMap
}
