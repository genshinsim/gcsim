package kinich

import "testing"

func TestBlindSpot(t *testing.T) {
	kinich := char{}

	kinich.characterAngularPosition = 0.
	kinich.blindSpotAngularPosition = 270.
	cross, boundary := kinich.NextMoveIsInBlindSpot(-1)
	if !cross {
		t.Errorf("expecting cross")
	}
	if boundary != kinich.blindSpotAngularPosition+blindSpotBoundary {
		t.Errorf("%v != %v", boundary, kinich.blindSpotAngularPosition+blindSpotBoundary)
	}

	kinich.characterAngularPosition = 50.
	cross, boundary = kinich.NextMoveIsInBlindSpot(1)
	if cross {
		t.Errorf("not expecting cross")
	}
	if boundary != -1 {
		t.Errorf("not expecting boundary")
	}

	kinich.characterAngularPosition = 324
	kinich.blindSpotAngularPosition = 0
	cross, boundary = kinich.NextMoveIsInBlindSpot(1)
	if !cross {
		t.Errorf("expecting cross")
	}
	if boundary != NormalizeAngle360(kinich.blindSpotAngularPosition-blindSpotBoundary) {
		t.Errorf("%v != %v", boundary, NormalizeAngle360(kinich.blindSpotAngularPosition-blindSpotBoundary))
	}

	kinich.characterAngularPosition = 216
	kinich.blindSpotAngularPosition = 180
	cross, boundary = kinich.NextMoveIsInBlindSpot(-1)
	if !cross {
		t.Errorf("expecting cross")
	}
	if boundary != kinich.blindSpotAngularPosition+blindSpotBoundary {
		t.Errorf("%v != %v", boundary, kinich.blindSpotAngularPosition+blindSpotBoundary)
	}

	kinich.characterAngularPosition = 0.
	kinich.blindSpotAngularPosition = 90.
	cross, boundary = kinich.NextMoveIsInBlindSpot(-1)
	if cross {
		t.Errorf("not expecting cross")
	}
	if boundary != -1 {
		t.Errorf("not expecting boundary")
	}
}
