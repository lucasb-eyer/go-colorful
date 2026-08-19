package colorful

import (
	"math"
	"testing"
)

// Reference values were produced with culori (https://culorijs.org), which
// implements Björn Ottosson's original Okhsl/Okhsv formulas.
var okhslVals = []struct {
	c     Color
	okhsl [3]float64
	okhsv [3]float64
}{
	{Color{1, 0, 0}, [3]float64{29.233880, 1.0, 0.568084}, [3]float64{29.233880, 0.999521, 1.0}},
	{Color{1, 1, 1}, [3]float64{0, 0, 1.0}, [3]float64{0, 0, 1.0}},
	{Color{0, 0, 0}, [3]float64{0, 0, 0}, [3]float64{0, 0, 0}},
	{Color{0.2, 0.2, 0.2}, [3]float64{0, 0, 0.220995}, [3]float64{0, 0, 0.220995}},
}

func TestOkhsl(t *testing.T) {
	for i, tt := range okhslVals {
		h, s, l := tt.c.Okhsl()
		// Hue is meaningless for achromatic colors.
		if tt.okhsl[1] != 0 && !almosteq(h, tt.okhsl[0]) {
			t.Errorf("%d. %v.Okhsl() h = %v, want %v", i, tt.c, h, tt.okhsl[0])
		}
		if !almosteq_eps(s, tt.okhsl[1], 1e-4) {
			t.Errorf("%d. %v.Okhsl() s = %v, want %v", i, tt.c, s, tt.okhsl[1])
		}
		if !almosteq_eps(l, tt.okhsl[2], 1e-4) {
			t.Errorf("%d. %v.Okhsl() l = %v, want %v", i, tt.c, l, tt.okhsl[2])
		}
	}
}

func TestOkhsv(t *testing.T) {
	for i, tt := range okhslVals {
		h, s, v := tt.c.Okhsv()
		if tt.okhsv[1] != 0 && !almosteq(h, tt.okhsv[0]) {
			t.Errorf("%d. %v.Okhsv() h = %v, want %v", i, tt.c, h, tt.okhsv[0])
		}
		if !almosteq_eps(s, tt.okhsv[1], 1e-4) {
			t.Errorf("%d. %v.Okhsv() s = %v, want %v", i, tt.c, s, tt.okhsv[1])
		}
		if !almosteq_eps(v, tt.okhsv[2], 1e-4) {
			t.Errorf("%d. %v.Okhsv() v = %v, want %v", i, tt.c, v, tt.okhsv[2])
		}
	}
}

// Okhsl(0, 1, 1) is white and Okhsl(0, 1, 0) is black regardless of hue and
// saturation, matching the reference implementation.
func TestOkhslExtremeLightness(t *testing.T) {
	if c := Okhsl(120, 1, 1); !c.AlmostEqualRgb(Color{1, 1, 1}) {
		t.Errorf("Okhsl(120, 1, 1) = %v, want white", c)
	}
	if c := Okhsl(120, 1, 0); !c.AlmostEqualRgb(Color{0, 0, 0}) {
		t.Errorf("Okhsl(120, 1, 0) = %v, want black", c)
	}
}

func TestOkhslRoundtrip(t *testing.T) {
	for r := 0; r <= 8; r++ {
		for g := 0; g <= 8; g++ {
			for b := 0; b <= 8; b++ {
				want := Color{float64(r) / 8, float64(g) / 8, float64(b) / 8}

				h, s, l := want.Okhsl()
				got := Okhsl(h, s, l)
				if !want.AlmostEqualRgb(got) {
					t.Errorf("Okhsl roundtrip: %v -> (%v, %v, %v) -> %v", want, h, s, l, got)
				}

				h, sv, v := want.Okhsv()
				got = Okhsv(h, sv, v)
				if !want.AlmostEqualRgb(got) {
					t.Errorf("Okhsv roundtrip: %v -> (%v, %v, %v) -> %v", want, h, sv, v, got)
				}
			}
		}
	}
}

// The forward conversions must stay finite and roughly within the nominal
// [0..1] range for every representable sRGB color. Saturation can overshoot 1.0
// by a hair for colors right on the gamut boundary, since the maximum
// saturation is derived from a polynomial approximation, so a small slack is
// allowed on the upper bound.
func TestOkhslRange(t *testing.T) {
	const slack = 0.03
	for r := 0; r <= 16; r++ {
		for g := 0; g <= 16; g++ {
			for b := 0; b <= 16; b++ {
				c := Color{float64(r) / 16, float64(g) / 16, float64(b) / 16}

				_, s, l := c.Okhsl()
				if math.IsNaN(s) || math.IsNaN(l) {
					t.Errorf("%v.Okhsl() produced NaN: s=%v l=%v", c, s, l)
				}
				if s < -slack || s > 1+slack || l < -1e-6 || l > 1+1e-6 {
					t.Errorf("%v.Okhsl() out of range: s=%v l=%v", c, s, l)
				}

				_, sv, v := c.Okhsv()
				if math.IsNaN(sv) || math.IsNaN(v) {
					t.Errorf("%v.Okhsv() produced NaN: s=%v v=%v", c, sv, v)
				}
				if sv < -slack || sv > 1+slack || v < -1e-6 || v > 1+1e-6 {
					t.Errorf("%v.Okhsv() out of range: s=%v v=%v", c, sv, v)
				}
			}
		}
	}
}
