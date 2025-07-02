# IP Tracker

A terminal-based IP geolocation tracker with an interactive world map built with Go and tview.

## Purpose

This project is created for educational purposes to study and learn the [tview](https://github.com/rivo/tview) library for building terminal user interfaces in Go. It demonstrates various tview components including custom widgets, mouse interactions, and dynamic color rendering.

## Features

- 🌍 **Interactive World Map**: Drag-and-drop world map with zoom and pan capabilities
- 📍 **IP Geolocation**: Automatically fetches and displays your current IP location
- 🎨 **Colored Markers**: Add custom colored markers to the map using string-based colors
- 🖱️ **Mouse Support**: Full mouse support for dragging and scrolling
- ⌨️ **Keyboard Navigation**: Arrow keys and keyboard shortcuts for navigation
- 🎭 **Dynamic Colors**: Beautiful color-coded information display

## Screenshots

The application displays:
- An interactive ASCII world map that you can drag and navigate
- Your current IP geolocation information with color-coded details
- Customizable markers on the map showing geographical positions

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd ip-tracker

# Build the application
go build -o ip-tracker.exe ./cmd/main.go

# Run the application
./ip-tracker.exe
```

## Usage

### Controls

- **Mouse**: Click and drag to pan the world map
- **Arrow Keys**: Navigate the map using keyboard
- **Home Key**: Return to map origin (0,0)
- **Mouse Wheel**: Scroll to zoom vertically
- **Ctrl+C**: Exit the application

### Adding Custom Markers

The WorldMapWidget supports adding custom markers programmatically:

```go
// Create a new world map widget
worldMap := components.NewWorldMapWidget()

// Add a point with latitude, longitude, color (as string), and symbol
worldMap.AddPoint(40.7128, -74.0060, "red", '●')    // New York
worldMap.AddPoint(51.5074, -0.1278, "blue", '★')    // London
worldMap.AddPoint(35.6762, 139.6503, "green", '▲')  // Tokyo
```

Supported colors include: `red`, `blue`, `green`, `yellow`, `cyan`, `magenta`, `white`, `lightblue`, `lightgreen`, `orange`, and more.

## Project Structure

```
ip-tracker/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── app/
│   │   └── ip_tracker.go       # Main application logic
│   ├── components/
│   │   └── worldmap.go         # Custom world map widget
│   ├── models/
│   │   └── ip_cloudflare.go    # Data models
│   └── services/
│       ├── ip_cloudflare.go    # Cloudflare API service
│       └── map.go              # Map utilities
├── pkg/
│   └── utils/
│       └── math.go             # Utility functions
└── assets/
    ├── map.txt                 # ASCII world map data
    └── map2.txt                # Alternative map data
```

## Technical Details

### Architecture

- **MVC Pattern**: Separation of concerns with models, services, and components
- **Custom Widgets**: Extended tview.Box to create a custom draggable map widget
- **String-based Colors**: Uses string color names that are converted to tcell colors internally
- **Coordinate Mapping**: Converts latitude/longitude to ASCII map positions

### Key Components

- **WorldMapWidget**: Custom tview widget supporting drag-and-drop with colored markers
- **MapPoint**: Structure representing geographical points with color and symbol
- **Cloudflare Service**: Fetches geolocation data from Cloudflare's API

## Acknowledgments

### APIs and Services

- **[Cloudflare Speed Test API](https://speed.cloudflare.com)**: This project uses Cloudflare's speed test metadata endpoint (`https://speed.cloudflare.com/meta`) to fetch IP geolocation information. We thank Cloudflare for providing this free service that enables accurate IP geolocation data.

### Libraries and Frameworks

- **[tview](https://github.com/rivo/tview)**: A rich interactive widget library for terminal applications in Go. This project serves as a study case for learning tview's capabilities, including custom widgets, mouse interactions, and dynamic color rendering. Special thanks to the tview maintainers for creating such a powerful and well-documented library.

- **[tcell](https://github.com/gdamore/tcell)**: Low-level terminal handling library that powers tview's rendering capabilities.

## Learning Objectives

This project demonstrates the following tview concepts:

1. **Custom Widget Development**: Creating a custom widget by extending `tview.Box`
2. **Mouse Event Handling**: Implementing drag-and-drop functionality
3. **Color Management**: Using both string-based and tcell color systems
4. **Layout Management**: Complex flex layouts with borders and titles
5. **Event Handling**: Keyboard and mouse input processing
6. **Dynamic Content**: Real-time updates and color-coded information display

## License

This project is open source and available under the [MIT License](LICENSE).

## Contributing

This is primarily an educational project for studying tview. However, contributions that enhance the learning aspects or add interesting tview features are welcome!

---

*Built with ❤️ using Go and tview for educational purposes*