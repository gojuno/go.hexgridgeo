package hexgridgeo

import (
	"math"
	"testing"

	hexgrid "github.com/gojuno/go.hexgrid"
)

func validatePoint(t *testing.T, e hexgrid.Point, r hexgrid.Point, precision float64) {
	if math.Abs(e.X()-r.X()) > precision || math.Abs(e.Y()-r.Y()) > precision {
		t.Errorf("expected point{x: %f, y: %f} but got point{x: %f, y: %f}", e.X(), e.Y(), r.X(), r.Y())
	}
}

func validateGeoPoint(t *testing.T, e Point, r Point, precision float64) {
	if math.Abs(e.Lon()-r.Lon()) > precision || math.Abs(e.Lat()-r.Lat()) > precision {
		t.Errorf("expected point{lon: %f, lat: %f} but got point{lon: %f, lat: %f}", e.Lon(), e.Lat(), r.Lon(), r.Lat())
	}
}

func TestProjectionNoOp(t *testing.T) {
	geoPoint := MakePoint(-73.0, 40.0)
	point := ProjectionNoOp.GeoToPoint(geoPoint)
	validatePoint(t, hexgrid.MakePoint(-73.0, 40.0), point, 0.00001)
	recodedGeoPoint := ProjectionNoOp.PointToGeo(point)
	validateGeoPoint(t, geoPoint, recodedGeoPoint, 0.00001)
}

func TestProjectionSin(t *testing.T) {
	geoPoint := MakePoint(-73.0, 40.0)
	point := ProjectionSin.GeoToPoint(geoPoint)
	validatePoint(t, hexgrid.MakePoint(9124497.47463, 4452779.63173), point, 0.00001)
	recodedGeoPoint := ProjectionSin.PointToGeo(point)
	validateGeoPoint(t, geoPoint, recodedGeoPoint, 0.00001)
}

func TestProjectionAEP(t *testing.T) {
	geoPoint := MakePoint(-73.0, 40.0)
	point := ProjectionAEP.GeoToPoint(geoPoint)
	validatePoint(t, hexgrid.MakePoint(-0.83453, -0.25514), point, 0.00001)
	recodedGeoPoint := ProjectionAEP.PointToGeo(point)
	validateGeoPoint(t, geoPoint, recodedGeoPoint, 0.00001)
}

func TestProjectionSM(t *testing.T) {
	geoPoint := MakePoint(-73.0, 40.0)
	point := ProjectionSM.GeoToPoint(geoPoint)
	validatePoint(t, hexgrid.MakePoint(-8126322.82791, 4865942.27950), point, 0.00001)
	recodedGeoPoint := ProjectionSM.PointToGeo(point)
	validateGeoPoint(t, geoPoint, recodedGeoPoint, 0.00001)
}

func TestSimple(t *testing.T) {
	grid := MakeGrid(OrientationFlat, 500, ProjectionSM)
	corners := grid.HexCorners(grid.HexAt(MakePoint(-73.0, 40.0)))
	expectedCorners := []Point{
		MakePoint(-72.99485, 39.99877), MakePoint(-72.99710, 40.00175),
		MakePoint(-73.00159, 40.00175), MakePoint(-73.00384, 39.99877),
		MakePoint(-73.00159, 39.99579), MakePoint(-72.99710, 39.99579)}
	for i := 0; i < 6; i++ {
		validateGeoPoint(t, expectedCorners[i], corners[i], 0.00001)
	}
}
