package kinich

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	const threshold = 0.0001
	return math.Abs(a-b) <= threshold
}

func TestBlindSpot(t *testing.T) {
	kinich := char{}

	kinich.characterAngularPosition = 0.
	kinich.blindSpotAngularPosition = 270.
	cross, boundary := kinich.NextMoveIsInBlindSpot(-1)
	if !cross {
		t.Errorf("expecting cross")
	}
	if !almostEqual(boundary, kinich.blindSpotAngularPosition+blindSpotBoundary) {
		t.Errorf("%v != %v", boundary, kinich.blindSpotAngularPosition+blindSpotBoundary)
	}

	kinich.characterAngularPosition = 50.
	cross, boundary = kinich.NextMoveIsInBlindSpot(1)
	if cross {
		t.Errorf("not expecting cross")
	}
	if !almostEqual(boundary, -1) {
		t.Errorf("not expecting boundary")
	}

	kinich.characterAngularPosition = 324
	kinich.blindSpotAngularPosition = 0
	cross, boundary = kinich.NextMoveIsInBlindSpot(1)
	if !cross {
		t.Errorf("expecting cross")
	}
	if !almostEqual(boundary, NormalizeAngle360(kinich.blindSpotAngularPosition-blindSpotBoundary)) {
		t.Errorf("%v != %v", boundary, NormalizeAngle360(kinich.blindSpotAngularPosition-blindSpotBoundary))
	}

	kinich.characterAngularPosition = 216
	kinich.blindSpotAngularPosition = 180
	cross, boundary = kinich.NextMoveIsInBlindSpot(-1)
	if !cross {
		t.Errorf("expecting cross")
	}
	if !almostEqual(boundary, kinich.blindSpotAngularPosition+blindSpotBoundary) {
		t.Errorf("%v != %v", boundary, kinich.blindSpotAngularPosition+blindSpotBoundary)
	}

	kinich.characterAngularPosition = 0.
	kinich.blindSpotAngularPosition = 90.
	cross, boundary = kinich.NextMoveIsInBlindSpot(-1)
	if cross {
		t.Errorf("not expecting cross")
	}
	if !almostEqual(boundary, -1) {
		t.Errorf("not expecting boundary")
	}

	kinich.characterAngularPosition = 359
	kinich.blindSpotAngularPosition = 35
	cross, boundary = kinich.NextMoveIsInBlindSpot(1)
	if !cross {
		t.Errorf("expecting cross")
	}
	if !almostEqual(boundary, kinich.blindSpotAngularPosition-blindSpotBoundary) {
		t.Errorf("%v != %v", boundary, kinich.blindSpotAngularPosition-blindSpotBoundary)
	}
	cross, boundary = kinich.NextMoveIsInBlindSpot(-1)
	if cross {
		t.Errorf("not expecting cross")
	}
	if boundary != -1 {
		t.Errorf("not expecting boundary")
	}
}
