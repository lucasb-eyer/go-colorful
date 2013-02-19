package colorful

import (
    "strings"
    "testing"
)

const delta = 1.0/256.0
func almosteq(v1, v2 float64) bool {
    return (v1 - v2) < delta && (v1 - v2) > -delta
}

// Note: the XYZ, L*a*b*, etc. are using D65 white and D50 white if postfixed by "50".
// See http://www.brucelindbloom.com/index.html?ColorCalcHelp.html
// For d50 white, no "adaptation" is used in colorful
var vals = []struct{
    c     Color
    hsv   [3]float64
    hex   string
    xyz   [3]float64
    lab   [3]float64
    lab50 [3]float64
    luv   [3]float64
    luv50 [3]float64
}{
    {Color{1.0, 1.0, 1.0}, [3]float64{  0.0, 0.0, 1.0}, "#ffffff", [3]float64{0.95050,1.00000,1.08900}, [3]float64{1.00000, 0.00005,-0.00100}, [3]float64{1.00000,-0.02388,-0.19362}, [3]float64{1.00000, 0.00001,-0.00017}, [3]float64{0.96422, 1.00000, 0.82521}},
    {Color{0.5, 1.0, 1.0}, [3]float64{180.0, 0.5, 1.0}, "#80ffff", [3]float64{0.62712,0.83329,1.07387}, [3]float64{0.93158,-0.35225,-0.10875}, [3]float64{0.93139,-0.37410,-0.30166}, [3]float64{0.93158,-0.53768,-0.11614}, [3]float64{0.62148, 0.82512, 0.81426}},
    {Color{1.0, 0.5, 1.0}, [3]float64{300.0, 0.5, 1.0}, "#ff80ff", [3]float64{0.67009,0.43918,0.99553}, [3]float64{0.72174, 0.64949,-0.42092}, [3]float64{0.72089, 0.63042,-0.61003}, [3]float64{0.72174, 0.59862,-0.77409}, [3]float64{0.66157, 0.43656, 0.74889}},
    {Color{1.0, 1.0, 0.5}, [3]float64{ 60.0, 0.5, 1.0}, "#ffff80", [3]float64{0.80896,0.94339,0.34368}, [3]float64{0.97768,-0.16538, 0.59979}, [3]float64{0.97768,-0.18802, 0.46795}, [3]float64{0.97768, 0.05742, 0.79595}, [3]float64{0.97768,-0.08646, 0.54509}},
    {Color{0.5, 0.5, 1.0}, [3]float64{240.0, 0.5, 1.0}, "#8080ff", [3]float64{0.34671,0.27248,0.98040}, [3]float64{0.59203, 0.33107,-0.63467}, [3]float64{0.59045, 0.31581,-0.82404}, [3]float64{0.59203,-0.07537,-1.04536}, [3]float64{0.31884, 0.26168, 0.73794}},
    {Color{1.0, 0.5, 0.5}, [3]float64{  0.0, 0.5, 1.0}, "#ff8080", [3]float64{0.52855,0.38257,0.25021}, [3]float64{0.68209, 0.48197, 0.22687}, [3]float64{0.68209, 0.46233, 0.10827}, [3]float64{0.68209, 0.91713, 0.19771}, [3]float64{0.68209, 0.81675, 0.02270}},
    {Color{0.5, 1.0, 0.5}, [3]float64{120.0, 0.5, 1.0}, "#80ff80", [3]float64{0.48558,0.77668,0.32854}, [3]float64{0.90628,-0.59893, 0.49697}, [3]float64{0.90628,-0.61803, 0.36710}, [3]float64{0.90628,-0.58686, 0.75862}, [3]float64{0.90628,-0.72023, 0.52608}},
    {Color{0.5, 0.5, 0.5}, [3]float64{  0.0, 0.0, 0.5}, "#808080", [3]float64{0.20518,0.21586,0.23507}, [3]float64{0.53585, 0.00003,-0.00006}, [3]float64{0.53585,-0.01429,-0.11622}, [3]float64{0.53585, 0.00000,-0.00009}, [3]float64{0.53585,-0.07885,-0.13758}},
    {Color{0.0, 1.0, 1.0}, [3]float64{180.0, 1.0, 1.0}, "#00ffff", [3]float64{0.53810,0.78740,1.06970}, [3]float64{0.91117,-0.48080,-0.14138}, [3]float64{0.91113,-0.50063,-0.33378}, [3]float64{0.91117,-0.70472,-0.15217}, [3]float64{0.52814, 0.77750, 0.81128}},
    {Color{1.0, 0.0, 1.0}, [3]float64{300.0, 1.0, 1.0}, "#ff00ff", [3]float64{0.59290,0.28480,0.96980}, [3]float64{0.60320, 0.98254,-0.60843}, [3]float64{0.60324, 0.96194,-0.79453}, [3]float64{0.60320, 0.84075,-1.08712}, [3]float64{0.96422, 1.00000, 0.82521}},
    {Color{1.0, 1.0, 0.0}, [3]float64{ 60.0, 1.0, 1.0}, "#ffff00", [3]float64{0.77000,0.92780,0.13850}, [3]float64{0.97138,-0.21556, 0.94482}, [3]float64{0.97138,-0.23782, 0.84745}, [3]float64{0.97138, 0.07703, 1.06789}, [3]float64{0.97138,-0.06592, 0.81865}},
    {Color{0.0, 0.0, 1.0}, [3]float64{240.0, 1.0, 1.0}, "#0000ff", [3]float64{0.18050,0.07220,0.95050}, [3]float64{0.32303, 0.79197,-1.07864}, [3]float64{0.32297, 0.77815,-1.26364}, [3]float64{0.32303,-0.09400,-1.30358}, [3]float64{0.14308, 0.06062, 0.71417}},
    {Color{0.0, 1.0, 0.0}, [3]float64{120.0, 1.0, 1.0}, "#00ff00", [3]float64{0.35760,0.71520,0.11920}, [3]float64{0.87737,-0.86185, 0.83181}, [3]float64{0.87737,-0.87909, 0.73919}, [3]float64{0.87737,-0.83080, 1.07401}, [3]float64{0.87737,-0.95991, 0.84890}},
    {Color{1.0, 0.0, 0.0}, [3]float64{  0.0, 1.0, 1.0}, "#ff0000", [3]float64{0.41240,0.21260,0.01930}, [3]float64{0.53233, 0.80109, 0.67220}, [3]float64{0.53233, 0.78301, 0.62172}, [3]float64{0.53233, 1.75053, 0.37751}, [3]float64{0.53233, 1.67219, 0.24092}},
    {Color{0.0, 0.0, 0.0}, [3]float64{  0.0, 0.0, 0.0}, "#000000", [3]float64{0.00000,0.00000,0.00000}, [3]float64{0.00000, 0.00000, 0.00000}, [3]float64{0.00000, 0.00000, 0.00000}, [3]float64{0.00000, 0.00000, 0.00000}, [3]float64{0.00000, 0.00000, 0.00000}},
}
var d50 = [3]float64{0.96421, 1.0, 0.82519}

/// HSV ///
///////////

func TestHsvCreation(t *testing.T) {
    for i, tt := range vals {
        c := Hsv(tt.hsv[0], tt.hsv[1], tt.hsv[2])
        if !c.AlmostEqual(tt.c) {
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

/// Hex ///
///////////

func TestHexCreation(t *testing.T) {
    for i, tt := range vals {
        c, err := Hex(tt.hex)
        if err != nil || !c.AlmostEqual(tt.c) {
            t.Errorf("%v. Hex(%v) => (%v), want %v (delta %v)", i, tt.hex, c, tt.c, delta)
        }
    }
}

func TestHEXCreation(t *testing.T) {
    for i, tt := range vals {
        c, err := Hex(strings.ToUpper(tt.hex))
        if err != nil || !c.AlmostEqual(tt.c) {
            t.Errorf("%v. HEX(%v) => (%v), want %v (delta %v)", i, strings.ToUpper(tt.hex), c, tt.c, delta)
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
// TODO

/// XYZ ///
///////////
func TestXyzCreation(t *testing.T) {
    for i, tt := range vals {
        c := Xyz(tt.xyz[0], tt.xyz[1], tt.xyz[2])
        if !c.AlmostEqual(tt.c) {
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


/// L*a*b* ///
//////////////
func TestLabCreation(t *testing.T) {
    for i, tt := range vals {
        c := Lab(tt.lab[0], tt.lab[1], tt.lab[2])
        if !c.AlmostEqual(tt.c) {
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
        c := LabWhiteRef(tt.lab50[0], tt.lab50[1], tt.lab50[2], d50[0], d50[1], d50[2])
        if !c.AlmostEqual(tt.c) {
            t.Errorf("%v. LabWhiteRef(%v, d50) => (%v), want %v (delta %v)", i, tt.lab50, c, tt.c, delta)
        }
    }
}

func TestLabWhiteRefConversion(t *testing.T) {
    for i, tt := range vals {
        l, a, b := tt.c.LabWhiteRef(d50[0], d50[1], d50[2])
        if !almosteq(l, tt.lab50[0]) || !almosteq(a, tt.lab50[1]) || !almosteq(b, tt.lab50[2]) {
            t.Errorf("%v. %v.LabWhiteRef(d50) => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, a, b}, tt.lab50, delta)
        }
    }
}

