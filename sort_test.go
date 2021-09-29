package colorful

import "testing"

// TestSortSimple tests the sorting of a small set of colors.
func TestSortSimple(t *testing.T) {
	// Sort a list of reds and blues.
	in := make([]Color, 0, 6)
	for i := 0; i < 3; i++ {
		in = append(in, Color{1.0 - float64(i+1)*0.25, 0.0, 0.0}) // Reds
		in = append(in, Color{0.0, 0.0, 1.0 - float64(i+1)*0.25}) // Blues
	}
	out := Sorted(in)

	// Ensure the output matches what we expected.
	exp := []Color{
		Color{R: 0.25, G: 0.0, B: 0},
		Color{R: 0.50, G: 0.0, B: 0},
		Color{R: 0.75, G: 0.0, B: 0},
		Color{R: 0.0, G: 0.0, B: 0.25},
		Color{R: 0.0, G: 0.0, B: 0.50},
		Color{R: 0.0, G: 0.0, B: 0.75},
	}
	for i, e := range exp {
		if out[i] != e {
			t.Fatalf("Expected %v but saw %v", e, out[i])
		}
	}
}
