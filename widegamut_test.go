package colorful

import (
	"math"
	"testing"
)

// Reference values computed using the CSS Color Level 4 spec's conversion code.
// Each entry maps an sRGB color to its representation in each wide-gamut space.
var widegamutVals = []struct {
	c        Color
	displayP3  [3]float64
	a98Rgb     [3]float64
	proPhoto   [3]float64
	rec2020    [3]float64
	xyzD50     [3]float64
}{
	// White
	{Color{1.0, 1.0, 1.0},
		[3]float64{1.0, 1.0, 1.0},
		[3]float64{1.0, 1.0, 1.0},
		[3]float64{1.0, 1.0, 1.0},
		[3]float64{1.0, 1.0, 1.0},
		[3]float64{0.9642200, 1.0000000, 0.8251900},
	},
	// Black
	{Color{0.0, 0.0, 0.0},
		[3]float64{0.0, 0.0, 0.0},
		[3]float64{0.0, 0.0, 0.0},
		[3]float64{0.0, 0.0, 0.0},
		[3]float64{0.0, 0.0, 0.0},
		[3]float64{0.0, 0.0, 0.0},
	},
}

/// Roundtrip tests ///
///////////////////////
// Encode sRGB -> space -> sRGB and verify identity within tolerance.

func TestDisplayP3Roundtrip(t *testing.T) {
	for _, tt := range vals {
		r, g, b := tt.c.DisplayP3()
		c := DisplayP3(r, g, b)
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v -> DisplayP3 -> sRGB = %v, want %v", tt.c, c, tt.c)
		}
	}
}

func TestA98RgbRoundtrip(t *testing.T) {
	for _, tt := range vals {
		r, g, b := tt.c.A98Rgb()
		c := A98Rgb(r, g, b)
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v -> A98Rgb -> sRGB = %v, want %v", tt.c, c, tt.c)
		}
	}
}

func TestProPhotoRgbRoundtrip(t *testing.T) {
	for _, tt := range vals {
		r, g, b := tt.c.ProPhotoRgb()
		c := ProPhotoRgb(r, g, b)
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v -> ProPhotoRgb -> sRGB = %v, want %v", tt.c, c, tt.c)
		}
	}
}

func TestRec2020Roundtrip(t *testing.T) {
	for _, tt := range vals {
		r, g, b := tt.c.Rec2020()
		c := Rec2020(r, g, b)
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v -> Rec2020 -> sRGB = %v, want %v", tt.c, c, tt.c)
		}
	}
}

func TestXyzD50Roundtrip(t *testing.T) {
	for _, tt := range vals {
		x, y, z := tt.c.XyzD50()
		c := XyzD50(x, y, z)
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v -> XyzD50 -> sRGB = %v, want %v", tt.c, c, tt.c)
		}
	}
}

/// Bradford adaptation roundtrip ///
/////////////////////////////////////

func TestBradfordRoundtrip(t *testing.T) {
	for _, tt := range vals {
		x, y, z := tt.c.Xyz()
		x50, y50, z50 := D65ToD50(x, y, z)
		x65, y65, z65 := D50ToD65(x50, y50, z50)
		if !almosteq(x, x65) || !almosteq(y, y65) || !almosteq(z, z65) {
			t.Errorf("Bradford roundtrip: (%v,%v,%v) -> D50 -> D65 = (%v,%v,%v)", x, y, z, x65, y65, z65)
		}
	}
}

/// Reference value tests ///
/////////////////////////////

func TestDisplayP3WhiteBlack(t *testing.T) {
	for i, tt := range widegamutVals {
		r, g, b := tt.c.DisplayP3()
		if !almosteq(r, tt.displayP3[0]) || !almosteq(g, tt.displayP3[1]) || !almosteq(b, tt.displayP3[2]) {
			t.Errorf("%v. %v.DisplayP3() => (%v, %v, %v), want %v", i, tt.c, r, g, b, tt.displayP3)
		}
	}
}

func TestA98RgbWhiteBlack(t *testing.T) {
	for i, tt := range widegamutVals {
		r, g, b := tt.c.A98Rgb()
		if !almosteq(r, tt.a98Rgb[0]) || !almosteq(g, tt.a98Rgb[1]) || !almosteq(b, tt.a98Rgb[2]) {
			t.Errorf("%v. %v.A98Rgb() => (%v, %v, %v), want %v", i, tt.c, r, g, b, tt.a98Rgb)
		}
	}
}

func TestProPhotoRgbWhiteBlack(t *testing.T) {
	for i, tt := range widegamutVals {
		r, g, b := tt.c.ProPhotoRgb()
		if !almosteq(r, tt.proPhoto[0]) || !almosteq(g, tt.proPhoto[1]) || !almosteq(b, tt.proPhoto[2]) {
			t.Errorf("%v. %v.ProPhotoRgb() => (%v, %v, %v), want %v", i, tt.c, r, g, b, tt.proPhoto)
		}
	}
}

func TestRec2020WhiteBlack(t *testing.T) {
	for i, tt := range widegamutVals {
		r, g, b := tt.c.Rec2020()
		if !almosteq(r, tt.rec2020[0]) || !almosteq(g, tt.rec2020[1]) || !almosteq(b, tt.rec2020[2]) {
			t.Errorf("%v. %v.Rec2020() => (%v, %v, %v), want %v", i, tt.c, r, g, b, tt.rec2020)
		}
	}
}

func TestXyzD50WhiteBlack(t *testing.T) {
	for i, tt := range widegamutVals {
		x, y, z := tt.c.XyzD50()
		if !almosteq(x, tt.xyzD50[0]) || !almosteq(y, tt.xyzD50[1]) || !almosteq(z, tt.xyzD50[2]) {
			t.Errorf("%v. %v.XyzD50() => (%v, %v, %v), want %v", i, tt.c, x, y, z, tt.xyzD50)
		}
	}
}

/// Transfer function tests ///
///////////////////////////////

func TestA98TransferRoundtrip(t *testing.T) {
	for v := 0.0; v <= 1.0; v += 0.01 {
		lin := linearizeA98(v)
		back := delinearizeA98(lin)
		if math.Abs(v-back) > 1e-10 {
			t.Errorf("A98 transfer roundtrip: %v -> %v -> %v", v, lin, back)
		}
	}
}

func TestProPhotoTransferRoundtrip(t *testing.T) {
	for v := 0.0; v <= 1.0; v += 0.01 {
		lin := linearizeProPhoto(v)
		back := delinearizeProPhoto(lin)
		if math.Abs(v-back) > 1e-10 {
			t.Errorf("ProPhoto transfer roundtrip: %v -> %v -> %v", v, lin, back)
		}
	}
}

func TestRec2020TransferRoundtrip(t *testing.T) {
	for v := 0.0; v <= 1.0; v += 0.01 {
		lin := linearizeRec2020(v)
		back := delinearizeRec2020(lin)
		if math.Abs(v-back) > 1e-10 {
			t.Errorf("Rec2020 transfer roundtrip: %v -> %v -> %v", v, lin, back)
		}
	}
}

/// sRGB primary roundtrip through each space ///
//////////////////////////////////////////////////
// sRGB red, green, blue through each wide-gamut space.

func TestSRGBPrimariesDisplayP3(t *testing.T) {
	primaries := []Color{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	for _, c := range primaries {
		r, g, b := c.DisplayP3()
		back := DisplayP3(r, g, b)
		if !back.AlmostEqualRgb(c) {
			t.Errorf("sRGB primary %v through DisplayP3 (%v,%v,%v) = %v", c, r, g, b, back)
		}
	}
}

func TestSRGBPrimariesA98Rgb(t *testing.T) {
	primaries := []Color{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	for _, c := range primaries {
		r, g, b := c.A98Rgb()
		back := A98Rgb(r, g, b)
		if !back.AlmostEqualRgb(c) {
			t.Errorf("sRGB primary %v through A98Rgb (%v,%v,%v) = %v", c, r, g, b, back)
		}
	}
}

func TestSRGBPrimariesProPhotoRgb(t *testing.T) {
	primaries := []Color{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	for _, c := range primaries {
		r, g, b := c.ProPhotoRgb()
		back := ProPhotoRgb(r, g, b)
		if !back.AlmostEqualRgb(c) {
			t.Errorf("sRGB primary %v through ProPhotoRgb (%v,%v,%v) = %v", c, r, g, b, back)
		}
	}
}

func TestSRGBPrimariesRec2020(t *testing.T) {
	primaries := []Color{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	for _, c := range primaries {
		r, g, b := c.Rec2020()
		back := Rec2020(r, g, b)
		if !back.AlmostEqualRgb(c) {
			t.Errorf("sRGB primary %v through Rec2020 (%v,%v,%v) = %v", c, r, g, b, back)
		}
	}
}

/// Blend endpoint tests ///
////////////////////////////

func TestBlendDisplayP3Endpoints(t *testing.T) {
	c1, _ := Hex("#1a1a46")
	c2, _ := Hex("#666666")
	if b := c1.BlendDisplayP3(c2, 0).Hex(); b != c1.Hex() {
		t.Errorf("BlendDisplayP3 t=0: got %v, want %v", b, c1.Hex())
	}
	if b := c1.BlendDisplayP3(c2, 1).Hex(); b != c2.Hex() {
		t.Errorf("BlendDisplayP3 t=1: got %v, want %v", b, c2.Hex())
	}
}

func TestBlendA98RgbEndpoints(t *testing.T) {
	c1, _ := Hex("#1a1a46")
	c2, _ := Hex("#666666")
	if b := c1.BlendA98Rgb(c2, 0).Hex(); b != c1.Hex() {
		t.Errorf("BlendA98Rgb t=0: got %v, want %v", b, c1.Hex())
	}
	if b := c1.BlendA98Rgb(c2, 1).Hex(); b != c2.Hex() {
		t.Errorf("BlendA98Rgb t=1: got %v, want %v", b, c2.Hex())
	}
}

func TestBlendProPhotoRgbEndpoints(t *testing.T) {
	c1, _ := Hex("#1a1a46")
	c2, _ := Hex("#666666")
	if b := c1.BlendProPhotoRgb(c2, 0).Hex(); b != c1.Hex() {
		t.Errorf("BlendProPhotoRgb t=0: got %v, want %v", b, c1.Hex())
	}
	if b := c1.BlendProPhotoRgb(c2, 1).Hex(); b != c2.Hex() {
		t.Errorf("BlendProPhotoRgb t=1: got %v, want %v", b, c2.Hex())
	}
}

func TestBlendRec2020Endpoints(t *testing.T) {
	c1, _ := Hex("#1a1a46")
	c2, _ := Hex("#666666")
	if b := c1.BlendRec2020(c2, 0).Hex(); b != c1.Hex() {
		t.Errorf("BlendRec2020 t=0: got %v, want %v", b, c1.Hex())
	}
	if b := c1.BlendRec2020(c2, 1).Hex(); b != c2.Hex() {
		t.Errorf("BlendRec2020 t=1: got %v, want %v", b, c2.Hex())
	}
}

/// XYZ D50 consistency with existing Lab D50 ///
/////////////////////////////////////////////////
// XyzD50 should produce values consistent with the existing D50 white point.

func TestXyzD50ConsistentWithLabD50(t *testing.T) {
	// Convert sRGB white to XYZ D50 and verify it matches the D50 white point
	x, y, z := Color{1, 1, 1}.XyzD50()
	if !almosteq(y, 1.0) {
		t.Errorf("White XYZ D50 Y = %v, want 1.0", y)
	}
	// D50 white point ratios should be close to the defined D50 constant
	if !almosteq_eps(x/y, D50[0]/D50[1], 0.005) {
		t.Errorf("White XYZ D50 x/y = %v, want ~%v", x/y, D50[0]/D50[1])
	}
	if !almosteq_eps(z/y, D50[2]/D50[1], 0.005) {
		t.Errorf("White XYZ D50 z/y = %v, want ~%v", z/y, D50[2]/D50[1])
	}
}
