package hexgridgeo

import (
	"math"

	hexgrid "github.com/gojuno/go.hexgrid"
	morton "github.com/gojuno/go.morton"
)

type Point struct {
	lon float64
	lat float64
}

type Projection interface {
	GeoToPoint(geoPoint Point) hexgrid.Point
	PointToGeo(point hexgrid.Point) Point
}

type projectionNoOp struct {
}

type projectionSin struct {
}

type projectionAEP struct {
}

type projectionSM struct {
}

var ProjectionNoOp = projectionNoOp{}
var ProjectionSin = projectionSin{}
var ProjectionAEP = projectionAEP{}
var ProjectionSM = projectionSM{}

var OrientationPointy = hexgrid.OrientationPointy
var OrientationFlat = hexgrid.OrientationFlat

type Grid struct {
	hexgrid    *hexgrid.Grid
	projection Projection
}

const earthCircumference = 40075016.685578488
const earthMetersPerDegree = 111319.49079327358

func MakePoint(lon float64, lat float64) Point {
	return Point{lon: lon, lat: lat}
}

func (point Point) Lon() float64 {
	return point.lon
}

func (point Point) Lat() float64 {
	return point.lat
}

func (projectionNoOp) GeoToPoint(geoPoint Point) hexgrid.Point {
	return hexgrid.MakePoint(geoPoint.Lon(), geoPoint.Lat())
}

func (projectionNoOp) PointToGeo(point hexgrid.Point) Point {
	return MakePoint(point.X(), point.Y())
}

func (projectionSin) GeoToPoint(geoPoint Point) hexgrid.Point {
	λ := (geoPoint.Lon() + 180) * (math.Pi / 180)
	φ := geoPoint.Lat() * (math.Pi / 180)
	x := (λ * math.Cos(φ)) * ((earthCircumference / 2) / math.Pi)
	y := φ * ((earthCircumference / 2) / math.Pi)
	return hexgrid.MakePoint(x, y)
}

func (projectionSin) PointToGeo(point hexgrid.Point) Point {
	φ := point.Y() / ((earthCircumference / 2) / math.Pi)
	λ := point.X() / (math.Cos(φ) * ((earthCircumference / 2) / math.Pi))
	lon := (λ / (math.Pi / 180)) - 180
	lat := φ / (math.Pi / 180)
	return MakePoint(lon, lat)
}

func (projectionAEP) GeoToPoint(geoPoint Point) hexgrid.Point {
	θ := geoPoint.Lon() * (math.Pi / 180)
	ρ := math.Pi/2 - (geoPoint.Lat() * (math.Pi / 180))
	x := ρ * math.Sin(θ)
	y := -ρ * math.Cos(θ)
	return hexgrid.MakePoint(x, y)
}

func (projectionAEP) PointToGeo(point hexgrid.Point) Point {
	θ := math.Atan2(point.X(), -point.Y())
	ρ := point.X() / math.Sin(θ)
	lat := (math.Pi/2 - ρ) / (math.Pi / 180)
	lon := θ / (math.Pi / 180)
	return MakePoint(lon, lat)
}

func (projectionSM) GeoToPoint(geoPoint Point) hexgrid.Point {
	latR := geoPoint.Lat() * (math.Pi / 180)
	x := geoPoint.Lon() * earthMetersPerDegree
	y := math.Log(math.Tan(latR) + (1 / math.Cos(latR)))
	y = (y / math.Pi) * (earthCircumference / 2)
	return hexgrid.MakePoint(x, y)
}

func (projectionSM) PointToGeo(point hexgrid.Point) Point {
	lon := point.X() / earthMetersPerDegree
	lat := math.Asin(math.Tanh((point.Y() / (earthCircumference / 2)) * math.Pi))
	lat = lat * (180 / math.Pi)
	return MakePoint(lon, lat)
}

func MakeGrid(orientation hexgrid.Orientation, size float64, projection Projection) *Grid {
	return &Grid{
		hexgrid:    hexgrid.MakeGrid(orientation, hexgrid.MakePoint(0, 0), hexgrid.MakePoint(size, size), morton.Make64(2, 32)),
		projection: projection}
}

func (grid *Grid) HexToCode(hex hexgrid.Hex) int64 {
	return grid.hexgrid.HexToCode(hex)
}

func (grid *Grid) HexFromCode(code int64) hexgrid.Hex {
	return grid.hexgrid.HexFromCode(code)
}

func (grid *Grid) HexAt(geoPoint Point) hexgrid.Hex {
	return grid.hexgrid.HexAt(grid.projection.GeoToPoint(geoPoint))
}

func (grid *Grid) HexCenter(hex hexgrid.Hex) Point {
	return grid.projection.PointToGeo(grid.hexgrid.HexCenter(hex))
}

func (grid *Grid) HexCorners(hex hexgrid.Hex) [6]Point {
	var geoCorners [6]Point
	corners := grid.hexgrid.HexCorners(hex)
	for i := 0; i < 6; i++ {
		geoCorners[i] = grid.projection.PointToGeo(corners[i])
	}
	return geoCorners
}

func (grid *Grid) HexNeighbors(hex hexgrid.Hex, layers int64) []hexgrid.Hex {
	return grid.hexgrid.HexNeighbors(hex, layers)
}

func (grid *Grid) MakeRegion(geometry []Point) *hexgrid.Region {
	points := make([]hexgrid.Point, len(geometry))
	for i := 0; i < len(geometry); i++ {
		points[i] = grid.projection.GeoToPoint(geometry[i])
	}
	return grid.hexgrid.MakeRegion(points)
}
