package colorful

import (
    "math"
    "strings"
    "testing"
)

// Checks whether the relative error is below the 8bit RGB delta, which should be good enough.
const delta = 1.0/256.0
func almosteq(v1, v2 float64) bool {
    if math.Abs(v1) > delta {
        return math.Abs((v1 - v2)/v1) < delta
    }
    return true
}

// Note: the XYZ, L*a*b*, etc. are using D65 white and D50 white if postfixed by "50".
// See http://www.brucelindbloom.com/index.html?ColorCalcHelp.html
// For d50 white, no "adaptation" and the sRGB model are used in colorful
// HCL values form http://www.easyrgb.com/index.php?X=CALC and missing ones hand-computed from lab ones
var vals = []struct{
    c     Color
    hsl   [3]float64
    hsv   [3]float64
    hex   string
    xyz   [3]float64
    xyy   [3]float64
    lab   [3]float64
    lab50 [3]float64
    luv   [3]float64
    luv50 [3]float64
    hcl   [3]float64
    hcl50 [3]float64
}{
    {Color{1.0, 1.0, 1.0}, [3]float64{  0.0, 0.0, 1.00}, [3]float64{  0.0, 0.0, 1.0}, "#ffffff", [3]float64{0.950470, 1.000000, 1.088830}, [3]float64{0.312727, 0.329023, 1.000000}, [3]float64{1.000000, 0.000000, 0.000000}, [3]float64{1.000000,-0.023881,-0.193622}, [3]float64{1.00000, 0.00000, 0.00000}, [3]float64{1.00000,-0.14716,-0.25658}, [3]float64{  0.0000, 0.000000, 1.000000}, [3]float64{262.9688, 0.195089, 1.000000}},
    {Color{0.5, 1.0, 1.0}, [3]float64{180.0, 1.0, 0.75}, [3]float64{180.0, 0.5, 1.0}, "#80ffff", [3]float64{0.626296, 0.832848, 1.073634}, [3]float64{0.247276, 0.328828, 0.832848}, [3]float64{0.931390,-0.353319,-0.108946}, [3]float64{0.931390,-0.374100,-0.301663}, [3]float64{0.93139,-0.53909,-0.11630}, [3]float64{0.93139,-0.67615,-0.35528}, [3]float64{197.1371, 0.369735, 0.931390}, [3]float64{218.8817, 0.480574, 0.931390}},
    {Color{1.0, 0.5, 1.0}, [3]float64{300.0, 1.0, 0.75}, [3]float64{300.0, 0.5, 1.0}, "#ff80ff", [3]float64{0.669430, 0.437920, 0.995150}, [3]float64{0.318397, 0.208285, 0.437920}, [3]float64{0.720892, 0.651673,-0.422133}, [3]float64{0.720892, 0.630425,-0.610035}, [3]float64{0.72089, 0.60047,-0.77626}, [3]float64{0.72089, 0.49438,-0.96123}, [3]float64{327.0661, 0.776450, 0.720892}, [3]float64{315.9417, 0.877257, 0.720892}},
    {Color{1.0, 1.0, 0.5}, [3]float64{ 60.0, 1.0, 0.75}, [3]float64{ 60.0, 0.5, 1.0}, "#ffff80", [3]float64{0.808654, 0.943273, 0.341930}, [3]float64{0.386203, 0.450496, 0.943273}, [3]float64{0.977637,-0.165795, 0.602017}, [3]float64{0.977637,-0.188424, 0.470410}, [3]float64{0.97764, 0.05759, 0.79816}, [3]float64{0.97764,-0.08628, 0.54731}, [3]float64{105.3975, 0.624430, 0.977637}, [3]float64{111.8287, 0.506743, 0.977637}},
    {Color{0.5, 0.5, 1.0}, [3]float64{240.0, 1.0, 0.75}, [3]float64{240.0, 0.5, 1.0}, "#8080ff", [3]float64{0.345256, 0.270768, 0.979954}, [3]float64{0.216329, 0.169656, 0.270768}, [3]float64{0.590453, 0.332846,-0.637099}, [3]float64{0.590453, 0.315806,-0.824040}, [3]float64{0.59045,-0.07568,-1.04877}, [3]float64{0.59045,-0.16257,-1.20027}, [3]float64{297.5843, 0.718805, 0.590453}, [3]float64{290.9689, 0.882482, 0.590453}},
    {Color{1.0, 0.5, 0.5}, [3]float64{  0.0, 1.0, 0.75}, [3]float64{  0.0, 0.5, 1.0}, "#ff8080", [3]float64{0.527613, 0.381193, 0.248250}, [3]float64{0.455996, 0.329451, 0.381193}, [3]float64{0.681085, 0.483884, 0.228328}, [3]float64{0.681085, 0.464258, 0.110043}, [3]float64{0.68108, 0.92148, 0.19879}, [3]float64{0.68108, 0.82125, 0.02404}, [3]float64{ 25.2610, 0.535049, 0.681085}, [3]float64{ 13.3347, 0.477121, 0.681085}},
    {Color{0.5, 1.0, 0.5}, [3]float64{120.0, 1.0, 0.75}, [3]float64{120.0, 0.5, 1.0}, "#80ff80", [3]float64{0.484480, 0.776121, 0.326734}, [3]float64{0.305216, 0.488946, 0.776121}, [3]float64{0.906026,-0.600870, 0.498993}, [3]float64{0.906026,-0.619946, 0.369365}, [3]float64{0.90603,-0.58869, 0.76102}, [3]float64{0.90603,-0.72202, 0.52855}, [3]float64{140.2920, 0.781050, 0.906026}, [3]float64{149.2134, 0.721640, 0.906026}},
    {Color{0.5, 0.5, 0.5}, [3]float64{  0.0, 0.0, 0.50}, [3]float64{  0.0, 0.0, 0.5}, "#808080", [3]float64{0.203440, 0.214041, 0.233054}, [3]float64{0.312727, 0.329023, 0.214041}, [3]float64{0.533890, 0.000000, 0.000000}, [3]float64{0.533890,-0.014285,-0.115821}, [3]float64{0.53389, 0.00000, 0.00000}, [3]float64{0.53389,-0.07857,-0.13699}, [3]float64{  0.0000, 0.000000, 0.533890}, [3]float64{262.9688, 0.116699, 0.533890}},
    {Color{0.0, 1.0, 1.0}, [3]float64{180.0, 1.0, 0.50}, [3]float64{180.0, 1.0, 1.0}, "#00ffff", [3]float64{0.538014, 0.787327, 1.069496}, [3]float64{0.224656, 0.328760, 0.787327}, [3]float64{0.911132,-0.480875,-0.141312}, [3]float64{0.911132,-0.500630,-0.333781}, [3]float64{0.91113,-0.70477,-0.15204}, [3]float64{0.91113,-0.83886,-0.38582}, [3]float64{196.3762, 0.501209, 0.911132}, [3]float64{213.6923, 0.601698, 0.911132}},
    {Color{1.0, 0.0, 1.0}, [3]float64{300.0, 1.0, 0.50}, [3]float64{300.0, 1.0, 1.0}, "#ff00ff", [3]float64{0.592894, 0.284848, 0.969638}, [3]float64{0.320938, 0.154190, 0.284848}, [3]float64{0.603242, 0.982343,-0.608249}, [3]float64{0.603242, 0.961939,-0.794531}, [3]float64{0.60324, 0.84071,-1.08683}, [3]float64{0.60324, 0.75194,-1.24161}, [3]float64{328.2350, 1.155407, 0.603242}, [3]float64{320.4444, 1.247640, 0.603242}},
    {Color{1.0, 1.0, 0.0}, [3]float64{ 60.0, 1.0, 0.50}, [3]float64{ 60.0, 1.0, 1.0}, "#ffff00", [3]float64{0.770033, 0.927825, 0.138526}, [3]float64{0.419320, 0.505246, 0.927825}, [3]float64{0.971393,-0.215537, 0.944780}, [3]float64{0.971393,-0.237800, 0.847398}, [3]float64{0.97139, 0.07706, 1.06787}, [3]float64{0.97139,-0.06590, 0.81862}, [3]float64{102.8512, 0.969054, 0.971393}, [3]float64{105.6754, 0.880131, 0.971393}},
    {Color{0.0, 0.0, 1.0}, [3]float64{240.0, 1.0, 0.50}, [3]float64{240.0, 1.0, 1.0}, "#0000ff", [3]float64{0.180437, 0.072175, 0.950304}, [3]float64{0.150000, 0.060000, 0.072175}, [3]float64{0.322970, 0.791875,-1.078602}, [3]float64{0.322970, 0.778150,-1.263638}, [3]float64{0.32297,-0.09405,-1.30342}, [3]float64{0.32297,-0.14158,-1.38629}, [3]float64{306.2849, 1.338076, 0.322970}, [3]float64{301.6248, 1.484014, 0.322970}},
    {Color{0.0, 1.0, 0.0}, [3]float64{120.0, 1.0, 0.50}, [3]float64{120.0, 1.0, 1.0}, "#00ff00", [3]float64{0.357576, 0.715152, 0.119192}, [3]float64{0.300000, 0.600000, 0.715152}, [3]float64{0.877347,-0.861827, 0.831793}, [3]float64{0.877347,-0.879067, 0.739170}, [3]float64{0.87735,-0.83078, 1.07398}, [3]float64{0.87735,-0.95989, 0.84887}, [3]float64{136.0160, 1.197759, 0.877347}, [3]float64{139.9409, 1.148534, 0.877347}},
    {Color{1.0, 0.0, 0.0}, [3]float64{  0.0, 1.0, 0.50}, [3]float64{  0.0, 1.0, 1.0}, "#ff0000", [3]float64{0.412456, 0.212673, 0.019334}, [3]float64{0.640000, 0.330000, 0.212673}, [3]float64{0.532408, 0.800925, 0.672032}, [3]float64{0.532408, 0.782845, 0.621518}, [3]float64{0.53241, 1.75015, 0.37756}, [3]float64{0.53241, 1.67180, 0.24096}, [3]float64{ 39.9990, 1.045518, 0.532408}, [3]float64{ 38.4469, 0.999566, 0.532408}},
    {Color{0.0, 0.0, 0.0}, [3]float64{  0.0, 0.0, 0.00}, [3]float64{  0.0, 0.0, 0.0}, "#000000", [3]float64{0.000000, 0.000000, 0.000000}, [3]float64{0.312727, 0.329023, 0.000000}, [3]float64{0.000000, 0.000000, 0.000000}, [3]float64{0.000000, 0.000000, 0.000000}, [3]float64{0.00000, 0.00000, 0.00000}, [3]float64{0.00000, 0.00000, 0.00000}, [3]float64{  0.0000, 0.000000, 0.000000}, [3]float64{  0.0000, 0.000000, 0.000000}},
}

// For testing short-hex values, since the above contains '80' values.
var shorthexvals = []struct{
    c     Color
    hex   string
}{
    {Color{1.0, 1.0, 1.0}, "#fff"},
    {Color{0.6, 1.0, 1.0}, "#9ff"},
    {Color{1.0, 0.6, 1.0}, "#f9f"},
    {Color{1.0, 1.0, 0.6}, "#ff9"},
    {Color{0.6, 0.6, 1.0}, "#99f"},
    {Color{1.0, 0.6, 0.6}, "#f99"},
    {Color{0.6, 1.0, 0.6}, "#9f9"},
    {Color{0.6, 0.6, 0.6}, "#999"},
    {Color{0.0, 1.0, 1.0}, "#0ff"},
    {Color{1.0, 0.0, 1.0}, "#f0f"},
    {Color{1.0, 1.0, 0.0}, "#ff0"},
    {Color{0.0, 0.0, 1.0}, "#00f"},
    {Color{0.0, 1.0, 0.0}, "#0f0"},
    {Color{1.0, 0.0, 0.0}, "#f00"},
    {Color{0.0, 0.0, 0.0}, "#000"},
}

/// HSV ///
///////////

func TestHsvCreation(t *testing.T) {
    for i, tt := range vals {
        c := Hsv(tt.hsv[0], tt.hsv[1], tt.hsv[2])
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Hsv(%v) => (%v), want %v (delta %v)", i, tt.hsv, c, tt.c, delta)
        }
    }
}

func TestHsvConversion(t *testing.T) {
    for i, tt := range vals {
        h, s, v := tt.c.Hsv()
        if !almosteq(h, tt.hsv[0]) || !almosteq(s, tt.hsv[1]) || !almosteq(v, tt.hsv[2]) {
            t.Errorf("%v. %v.Hsv() => (%v), want %v (delta %v)", i, tt.c, []float64{h, s, v}, tt.hsv, delta)
        }
    }
}

/// HSL ///
///////////

func TestHslCreation(t *testing.T) {
    for i, tt := range vals {
        c := Hsl(tt.hsl[0], tt.hsl[1], tt.hsl[2])
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Hsl(%v) => (%v), want %v (delta %v)", i, tt.hsl, c, tt.c, delta)
        }
    }
}

func TestHslConversion(t *testing.T) {
    for i, tt := range vals {
        h, s, l := tt.c.Hsl()
        if !almosteq(h, tt.hsl[0]) || !almosteq(s, tt.hsl[1]) || !almosteq(l, tt.hsl[2]) {
            t.Errorf("%v. %v.Hsl() => (%v), want %v (delta %v)", i, tt.c, []float64{h, s, l}, tt.hsl, delta)
        }
    }
}

/// Hex ///
///////////

func TestHexCreation(t *testing.T) {
    for i, tt := range vals {
        c, err := Hex(tt.hex)
        if err != nil || !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Hex(%v) => (%v), want %v (delta %v)", i, tt.hex, c, tt.c, delta)
        }
    }
}

func TestHEXCreation(t *testing.T) {
    for i, tt := range vals {
        c, err := Hex(strings.ToUpper(tt.hex))
        if err != nil || !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. HEX(%v) => (%v), want %v (delta %v)", i, strings.ToUpper(tt.hex), c, tt.c, delta)
        }
    }
}

func TestShortHexCreation(t *testing.T) {
    for i, tt := range shorthexvals {
        c, err := Hex(tt.hex)
        if err != nil || !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Hex(%v) => (%v), want %v (delta %v)", i, tt.hex, c, tt.c, delta)
        }
    }
}

func TestShortHEXCreation(t *testing.T) {
    for i, tt := range shorthexvals {
        c, err := Hex(strings.ToUpper(tt.hex))
        if err != nil || !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Hex(%v) => (%v), want %v (delta %v)", i, strings.ToUpper(tt.hex), c, tt.c, delta)
        }
    }
}

func TestHexConversion(t *testing.T) {
    for i, tt := range vals {
        hex := tt.c.Hex()
        if hex != tt.hex {
            t.Errorf("%v. %v.Hex() => (%v), want %v (delta %v)", i, tt.c, hex, tt.hex, delta)
        }
    }
}

/// Linear ///
//////////////
// TODO (implicitly tested by XYZ)

/// XYZ ///
///////////
func TestXyzCreation(t *testing.T) {
    for i, tt := range vals {
        c := Xyz(tt.xyz[0], tt.xyz[1], tt.xyz[2])
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Xyz(%v) => (%v), want %v (delta %v)", i, tt.xyz, c, tt.c, delta)
        }
    }
}

func TestXyzConversion(t *testing.T) {
    for i, tt := range vals {
        x, y, z := tt.c.Xyz()
        if !almosteq(x, tt.xyz[0]) || !almosteq(y, tt.xyz[1]) || !almosteq(z, tt.xyz[2]) {
            t.Errorf("%v. %v.Xyz() => (%v), want %v (delta %v)", i, tt.c, [3]float64{x, y, z}, tt.xyz, delta)
        }
    }
}


/// xyY ///
///////////
func TestXyyCreation(t *testing.T) {
    for i, tt := range vals {
        c := Xyy(tt.xyy[0], tt.xyy[1], tt.xyy[2])
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Xyy(%v) => (%v), want %v (delta %v)", i, tt.xyy, c, tt.c, delta)
        }
    }
}

func TestXyyConversion(t *testing.T) {
    for i, tt := range vals {
        x, y, Y := tt.c.Xyy()
        if !almosteq(x, tt.xyy[0]) || !almosteq(y, tt.xyy[1]) || !almosteq(Y, tt.xyy[2]) {
            t.Errorf("%v. %v.Xyy() => (%v), want %v (delta %v)", i, tt.c, [3]float64{x, y, Y}, tt.xyy, delta)
        }
    }
}


/// L*a*b* ///
//////////////
func TestLabCreation(t *testing.T) {
    for i, tt := range vals {
        c := Lab(tt.lab[0], tt.lab[1], tt.lab[2])
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Lab(%v) => (%v), want %v (delta %v)", i, tt.lab, c, tt.c, delta)
        }
    }
}

func TestLabConversion(t *testing.T) {
    for i, tt := range vals {
        l, a, b := tt.c.Lab()
        if !almosteq(l, tt.lab[0]) || !almosteq(a, tt.lab[1]) || !almosteq(b, tt.lab[2]) {
            t.Errorf("%v. %v.Lab() => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, a, b}, tt.lab, delta)
        }
    }
}

func TestLabWhiteRefCreation(t *testing.T) {
    for i, tt := range vals {
        c := LabWhiteRef(tt.lab50[0], tt.lab50[1], tt.lab50[2], D50)
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. LabWhiteRef(%v, D50) => (%v), want %v (delta %v)", i, tt.lab50, c, tt.c, delta)
        }
    }
}

func TestLabWhiteRefConversion(t *testing.T) {
    for i, tt := range vals {
        l, a, b := tt.c.LabWhiteRef(D50)
        if !almosteq(l, tt.lab50[0]) || !almosteq(a, tt.lab50[1]) || !almosteq(b, tt.lab50[2]) {
            t.Errorf("%v. %v.LabWhiteRef(D50) => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, a, b}, tt.lab50, delta)
        }
    }
}


/// L*u*v* ///
//////////////
func TestLuvCreation(t *testing.T) {
    for i, tt := range vals {
        c := Luv(tt.luv[0], tt.luv[1], tt.luv[2])
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Luv(%v) => (%v), want %v (delta %v)", i, tt.luv, c, tt.c, delta)
        }
    }
}

func TestLuvConversion(t *testing.T) {
    for i, tt := range vals {
        l, u, v := tt.c.Luv()
        if !almosteq(l, tt.luv[0]) || !almosteq(u, tt.luv[1]) || !almosteq(v, tt.luv[2]) {
            t.Errorf("%v. %v.Luv() => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, u, v}, tt.luv, delta)
        }
    }
}

func TestLuvWhiteRefCreation(t *testing.T) {
    for i, tt := range vals {
        c := LuvWhiteRef(tt.luv50[0], tt.luv50[1], tt.luv50[2], D50)
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. LuvWhiteRef(%v, D50) => (%v), want %v (delta %v)", i, tt.luv50, c, tt.c, delta)
        }
    }
}

func TestLuvWhiteRefConversion(t *testing.T) {
    for i, tt := range vals {
        l, u, v := tt.c.LuvWhiteRef(D50)
        if !almosteq(l, tt.luv50[0]) || !almosteq(u, tt.luv50[1]) || !almosteq(v, tt.luv50[2]) {
            t.Errorf("%v. %v.LuvWhiteRef(D50) => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, u, v}, tt.luv50, delta)
        }
    }
}


/// HCL ///
///////////
// CIE-L*a*b* in polar coordinates.
func TestHclCreation(t *testing.T) {
    for i, tt := range vals {
        c := Hcl(tt.hcl[0], tt.hcl[1], tt.hcl[2])
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. Hcl(%v) => (%v), want %v (delta %v)", i, tt.hcl, c, tt.c, delta)
        }
    }
}

func TestHclConversion(t *testing.T) {
    for i, tt := range vals {
        h, c, l := tt.c.Hcl()
        if !almosteq(h, tt.hcl[0]) || !almosteq(c, tt.hcl[1]) || !almosteq(l, tt.hcl[2]) {
            t.Errorf("%v. %v.Hcl() => (%v), want %v (delta %v)", i, tt.c, [3]float64{h, c, l}, tt.hcl, delta)
        }
    }
}

func TestHclWhiteRefCreation(t *testing.T) {
    for i, tt := range vals {
        c := HclWhiteRef(tt.hcl50[0], tt.hcl50[1], tt.hcl50[2], D50)
        if !c.AlmostEqualRgb(tt.c) {
            t.Errorf("%v. HclWhiteRef(%v, D50) => (%v), want %v (delta %v)", i, tt.hcl50, c, tt.c, delta)
        }
    }
}

func TestHclWhiteRefConversion(t *testing.T) {
    for i, tt := range vals {
        h, c, l := tt.c.HclWhiteRef(D50)
        if !almosteq(h, tt.hcl50[0]) || !almosteq(c, tt.hcl50[1]) || !almosteq(l, tt.hcl50[2]) {
            t.Errorf("%v. %v.HclWhiteRef(D50) => (%v), want %v (delta %v)", i, tt.c, [3]float64{h, c, l}, tt.hcl50, delta)
        }
    }
}

