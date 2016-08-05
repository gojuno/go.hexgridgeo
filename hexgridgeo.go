package hexgridgeo

import (
	"math"

	hexgrid "github.com/gojuno/go.hexgrid"
	morton "github.com/gojuno/go.morton"
	geo "github.com/kellydunn/golang-geo"
)

type Projection interface {
	GeoToPoint(geoPoint *geo.Point) hexgrid.Point
	PointToGeo(point hexgrid.Point) *geo.Point
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

var PointyOrientation = hexgrid.PointyOrientation
var FlatOrientation = hexgrid.FlatOrientation

type Grid struct {
	hexgrid    *hexgrid.Grid
	projection Projection
}

const earthCircumference = 40075016.685578488
const earthMetersPerDegree = 111319.49079327358

func (projectionNoOp) GeoToPoint(geoPoint *geo.Point) hexgrid.Point {
	return hexgrid.MakePoint(geoPoint.Lng(), geoPoint.Lat())
}

func (projectionNoOp) PointToGeo(point hexgrid.Point) *geo.Point {
	return geo.NewPoint(point.Y(), point.X())
}

func (projectionSin) GeoToPoint(geoPoint *geo.Point) hexgrid.Point {
	λ := (geoPoint.Lng() + 180) * (math.Pi / 180)
	φ := geoPoint.Lat() * (math.Pi / 180)
	x := (λ * math.Cos(φ)) * ((earthCircumference / 2) / math.Pi)
	y := φ * ((earthCircumference / 2) / math.Pi)
	return hexgrid.MakePoint(x, y)
}

func (projectionSin) PointToGeo(point hexgrid.Point) *geo.Point {
	φ := point.Y() / ((earthCircumference / 2) / math.Pi)
	λ := point.X() / (math.Cos(φ) * ((earthCircumference / 2) / math.Pi))
	lon := (λ / (math.Pi / 180)) - 180
	lat := φ / (math.Pi / 180)
	return geo.NewPoint(lat, lon)
}

func (projectionAEP) GeoToPoint(geoPoint *geo.Point) hexgrid.Point {
	θ := geoPoint.Lng() * (math.Pi / 180)
	ρ := math.Pi/2 - (geoPoint.Lat() * (math.Pi / 180))
	x := ρ * math.Sin(θ)
	y := -ρ * math.Cos(θ)
	return hexgrid.MakePoint(x, y)
}

func (projectionAEP) PointToGeo(point hexgrid.Point) *geo.Point {
	θ := math.Atan2(point.X(), -point.Y())
	ρ := point.X() / math.Sin(θ)
	lat := (math.Pi/2 - ρ) / (math.Pi / 180)
	lon := θ / (math.Pi / 180)
	return geo.NewPoint(lat, lon)
}

func (projectionSM) GeoToPoint(geoPoint *geo.Point) hexgrid.Point {
	latR := geoPoint.Lat() * (math.Pi / 180)
	x := geoPoint.Lng() * earthMetersPerDegree
	y := math.Log(math.Tan(latR) + (1 / math.Cos(latR)))
	y = (y / math.Pi) * (earthCircumference / 2)
	return hexgrid.MakePoint(x, y)
}

func (projectionSM) PointToGeo(point hexgrid.Point) *geo.Point {
	lon := point.X() / earthMetersPerDegree
	lat := math.Asin(math.Tanh((point.Y() / (earthCircumference / 2)) * math.Pi))
	lat = lat * (180 / math.Pi)
	return geo.NewPoint(lat, lon)
}

func MakeGrid(orientation hexgrid.Orientation, size float64, projection Projection) *Grid {
	return &Grid{
		hexgrid:    hexgrid.MakeGrid(orientation, hexgrid.MakePoint(0, 0), hexgrid.MakePoint(size, size), morton.Make64(2, 31)),
		projection: projection}
}

func (grid *Grid) HexToCode(hex hexgrid.Hex) uint64 {
	return grid.hexgrid.HexToCode(hex)
}

func (grid *Grid) HexFromCode(code uint64) hexgrid.Hex {
	return grid.hexgrid.HexFromCode(code)
}

func (grid *Grid) HexAt(geoPoint *geo.Point) hexgrid.Hex {
	return grid.hexgrid.HexAt(grid.projection.GeoToPoint(geoPoint))
}

func (grid *Grid) HexCenter(hex hexgrid.Hex) *geo.Point {
	return grid.projection.PointToGeo(grid.hexgrid.HexCenter(hex))
}

func (grid *Grid) HexCorners(hex hexgrid.Hex) [6]*geo.Point {
	var geoCorners [6]*geo.Point
	corners := grid.hexgrid.HexCorners(hex)
	for i := 0; i < 6; i++ {
		geoCorners[i] = grid.projection.PointToGeo(corners[i])
	}
	return geoCorners
}

func (grid *Grid) HexNeighbors(hex hexgrid.Hex, layers int64) []hexgrid.Hex {
	return grid.hexgrid.HexNeighbors(hex, layers)
}

func (grid *Grid) MakeRegion(geometry *geo.Polygon) *hexgrid.Region {
	geoPoints := geometry.Points()
	points := make([]hexgrid.Point, len(geoPoints))

	for i := 0; i < len(geoPoints); i++ {
		points[i] = grid.projection.GeoToPoint(geoPoints[i])
	}

	return grid.hexgrid.MakeRegion(points)
}
