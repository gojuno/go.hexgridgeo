# HexGrid GEO [![GoDoc](https://godoc.org/github.com/gojuno/go.hexgridgeo?status.svg)](http://godoc.org/github.com/gojuno/go.hexgridgeo) [![Build Status](https://travis-ci.org/gojuno/go.hexgridgeo.svg?branch=master)](https://travis-ci.org/gojuno/go.hexgridgeo)

## Basics

GEO wrapper for [[https://github.com/gojuno/go.hexgrid][HexGrid]].

## Examples

```
import "github.com/gojuno/go.hexgridgeo"

grid := hexgridgeo.MakeGrid(hexgridgeo.OrientationFlat, 500, hexgridgeo.ProjectionSM)
hex := grid.HexAt(hexgridgeo.MakePoint(-73.5, 40.3))
code := grid.HexToCode(hex)
restoredHex := grid.HexFromCode(code)
neighbors := grid.HexNeighbors(hex, 2)
points := []hexgridgeo.Point{hexgridgeo.MakePoint(-73.0, 40.0), hexgridgeo.MakePoint(-74.0, 40.0), hexgridgeo.MakePoint(-74.0, 41.0), hexgridgeo.MakePoint(-73.0, 41.0)}
region := grid.MakeRegion(points)
hexesInRegion := region.Hexes()
```