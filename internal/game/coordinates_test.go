package game

import "testing"

// TestTileWorldCenterUsesSanctumOrigin verifies Tile centers use the Sanctum as origin.
func TestTileWorldCenterUsesSanctumOrigin(t *testing.T) {
	tests := []struct {
		name string
		x    int
		y    int
		want worldPosition
	}{
		{name: "sanctum", x: homePlotCenter, y: homePlotCenter, want: worldPosition{X: 0, Y: 0}},
		{name: "north road", x: homePlotCenter, y: homePlotCenter - 1, want: worldPosition{X: 0, Y: 1}},
		{name: "east", x: homePlotCenter + 1, y: homePlotCenter, want: worldPosition{X: 1, Y: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tileWorldCenter(tt.x, tt.y); got != tt.want {
				t.Fatalf("tileWorldCenter(%d, %d) = %+v, want %+v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

// TestNorthRoadSharedEdgeCoordinate verifies the documented north-road edge example.
func TestNorthRoadSharedEdgeCoordinate(t *testing.T) {
	_, northEdge, _, height := tileWorldRect(homePlotCenter, homePlotCenter-1)
	got := northEdge
	want := 1.5
	if got != want {
		t.Fatalf("first north road Tile north edge = %f, want shared edge %f", got, want)
	}
	if height != 1 {
		t.Fatalf("Tile height = %f, want 1", height)
	}
}
