package models

type MapPosition struct {
    X int // Column position in the ASCII map
    Y int // Row position in the ASCII map
}

// MapDimensions represents the ASCII map dimensions
type MapDimensions struct {
    Width  int
    Height int
}