// The colorful package provides all kinds of functions for working with colors.
package colorful

import(
    "fmt"
    "math"
    "math/rand"
)

// A color is stored internally using sRGB (standard RGB) values in the range 0-1
type Color struct {
    R, G, B float64
}

// This is the tolerance used when comparing colors using AlmostEqual.
const Delta = 1.0/255.0

// Check for equality between colors within the tolerance Delta (1/255).
func (c1 Color) AlmostEqual(c2 Color) bool {
    return math.Abs(c1.R - c2.R) +
           math.Abs(c1.G - c2.G) +
           math.Abs(c1.B - c2.B) < 3.0*Delta
}

// BlendLab blends two colors in the L*a*b* color-space, which should result in a smoother blend.
// t == 0 results in c1, t == 1 results in c2
func (c1 Color) BlendLab(c2 Color, t float64) Color {
    l1, a1, b1 := c1.Lab()
    l2, a2, b2 := c2.Lab()
    return Lab(l1 + t*(l2 - l1),
               a1 + t*(a2 - a1),
               b1 + t*(b2 - b1))
}

//func (c1 Color) BlendHcl(c2 Color, t float64) Color {
//}

// TODO: implement HCL, maybe drop HSV?
// http://stat.ethz.ch/R-manual/R-patched/library/grDevices/html/hcl.html
// http://stackoverflow.com/questions/7530627/hcl-color-to-rgb-and-backward

/// HSV ///
///////////
// From http://en.wikipedia.org/wiki/HSL_and_HSV
// Note that h is in [0..360] and s,v in [0..1]

// Hsv returns the Hue [0..360], Saturation and Value [0..1] of the color.
func (col Color) Hsv() (h, s, v float64) {
    min := math.Min(math.Min(col.R, col.G), col.B)
    v    = math.Max(math.Max(col.R, col.G), col.B)
    C := v - min

    s = 0.0
    if v != 0.0 {
        s = C / v
    }

    h = 0.0  // We use 0 instead of undefined as in wp.
    if min != v {
        if v == col.R { h = math.Mod((col.G - col.B) / C, 6.0) }
        if v == col.G { h = (col.B - col.R) / C + 2.0 }
        if v == col.B { h = (col.R - col.G) / C + 4.0 }
        h *= 60.0
        if h < 0.0 { h += 360.0 }
    }
    return h, s, v
}

// Hsv creates a new Color given a Hue in [0..360], a Saturation and a Value in [0..1]
func Hsv(H, S, V float64) Color {
    Hp := H/60.0
    C := V*S
    X := C*(1.0-math.Abs(math.Mod(Hp, 2.0)-1.0))

    m := V-C;
    r, g, b := 0.0, 0.0, 0.0

    switch {
    case 0.0 <= Hp && Hp < 1.0: r = C; g = X
    case 1.0 <= Hp && Hp < 2.0: r = X; g = C
    case 2.0 <= Hp && Hp < 3.0: g = C; b = X
    case 3.0 <= Hp && Hp < 4.0: g = X; b = C
    case 4.0 <= Hp && Hp < 5.0: r = X; b = C
    case 5.0 <= Hp && Hp < 6.0: r = C; b = X
    }

    return Color{m+r, m+g, m+b}
}

/// Hex ///
///////////

// Hex returns the hex "html" representation of the color, as in #ff0080.
func (col Color) Hex() string {
    // Add 0.5 for rounding
    return fmt.Sprintf("#%02x%02x%02x", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
}

// Hex parses a "html" hex color-string, either in the 3 "#f0c" or 6 "#ff1034" digits form.
func Hex(scol string) (col Color, err error) {
    format := "#%02x%02x%02x"
    if len(scol) == 4 {
        format = "#%x%x%x"
    }

    var r, g, b uint8
    n, err := fmt.Sscanf(scol, format, &r, &g, &b)
    if err != nil {
        return Color{}, err
    }
    if n != 3 {
        return Color{}, fmt.Errorf("color: %v is not a hex-color", scol)
    }

    return Color{float64(r)/255.0, float64(g)/255.0, float64(b)/255.0}, nil
}

/// Linear ///
//////////////
// http://www.sjbrown.co.uk/2004/05/14/gamma-correct-rendering/

func linearize(v float64) float64 {
    if v <= 0.04045 {
        return v / 12.92
    }
    return math.Pow((v + 0.055)/1.055, 2.4)
}

// LinearRgb converts the color into the linear RGB space (see http://www.sjbrown.co.uk/2004/05/14/gamma-correct-rendering/).
func (col Color) LinearRgb() (r, g, b float64) {
    r = linearize(col.R)
    g = linearize(col.G)
    b = linearize(col.B)
    return r, g, b
}

// FastLinearRgb is much faster than and almost as accurate as LinearRgb.
func (col Color) FastLinearRgb() (r, g, b float64) {
    r = math.Pow(col.R, 2.2)
    g = math.Pow(col.G, 2.2)
    b = math.Pow(col.B, 2.2)
    return r, g, b
}

func delinearize(v float64) float64 {
    if v <= 0.0031308 {
        return 12.92 * v
    }
    return 1.055 * math.Pow(v, 1.0/2.4) - 0.055
}

// LinearRgb creates an sRGB color out of the given linear RGB color (see http://www.sjbrown.co.uk/2004/05/14/gamma-correct-rendering/).
func LinearRgb(r, g, b float64) Color {
    return Color{delinearize(r), delinearize(g), delinearize(b)}
}

// FastLinearRgb is much faster than and almost as accurate as LinearRgb.
func FastLinearRgb(r, g, b float64) Color {
    return Color{math.Pow(r, 1.0/2.2), math.Pow(g, 1.0/2.2), math.Pow(b, 1.0/2.2)}
}

// XyzToLinearRgb converts from CIE XYZ-space to Linear RGB space.
func XyzToLinearRgb(x, y, z float64) (r, g, b float64) {
    r =  3.2406*x - 1.5372*y - 0.4986*z
    g = -0.9689*x + 1.8758*y + 0.0416*z
    b =  0.0557*x - 0.2040*y + 1.0570*z
    return r, g, b
}

func LinearRgbToXyz(r, g, b float64) (x, y, z float64) {
    x = 0.4124*r + 0.3576*g + 0.1805*b
    y = 0.2126*r + 0.7152*g + 0.0722*b
    z = 0.0193*r + 0.1192*g + 0.9505*b
    return x, y, z
}

/// XYZ ///
///////////
// http://www.sjbrown.co.uk/2004/05/14/gamma-correct-rendering/

func (col Color) Xyz() (float64, float64, float64) {
    return LinearRgbToXyz(col.LinearRgb())
}

func Xyz(x, y, z float64) Color {
    return LinearRgb(XyzToLinearRgb(x, y, z))
}

/// L*a*b* ///
//////////////
// http://en.wikipedia.org/wiki/Lab_color_space#CIELAB-CIEXYZ_conversions

var d65 = [3]float64{0.95043, 1.00000, 1.08890}

func lab_f(t float64) float64 {
    if t > 6.0/29.0 * 6.0/29.0 * 6.0/29.0 {
        return math.Cbrt(t)
    }
    return t/3.0 * 29.0/6.0 * 29.0/6.0 + 4.0/29.0
}

// For L*a*b*, we need to L*a*b*->XYZ->RGB and the first one is device dependent.
func XyzToLab(x, y, z float64) (l, a, b float64) {
    // Use D65 white as reference point by default.
    // http://www.fredmiranda.com/forum/topic/1035332
    // http://en.wikipedia.org/wiki/Standard_illuminant
    return XyzToLabWhiteRef(x, y, z, d65[0], d65[1], d65[2])
}

func XyzToLabWhiteRef(x, y, z, Xwref, Ywref, Zwref float64) (l, a, b float64) {
    fy := lab_f(y/Ywref)
    l = 1.16*fy - 0.16
    a = 5.0*(lab_f(x/Xwref) - fy)
    b = 2.0*(fy - lab_f(z/Zwref))
    return l, a, b
}

func lab_finv(t float64) float64 {
    if t > 6.0/29.0 {
        return t * t * t
    }
    return 3.0 * 6.0/29.0 * 6.0/29.0 * (t - 4.0/29.0)
}

func LabToXyz(l, a, b float64) (x, y, z float64) {
    // D65 white (see above).
    return LabToXyzWhiteRef(l, a, b, d65[0], d65[1], d65[2])
}

func LabToXyzWhiteRef(l, a, b, Xwref, Ywref, Zwref float64) (x, y, z float64) {
    l2 := (l + 0.16) / 1.16
    x = Xwref * lab_finv(l2 + a/5.0)
    y = Ywref * lab_finv(l2)
    z = Zwref * lab_finv(l2 - b/2.0)
    return x, y, z
}

func (col Color) Lab() (float64, float64, float64) {
    return XyzToLab(col.Xyz())
}

func (col Color) LabWhiteRef(Xwref, Ywref, Zwref float64) (float64, float64, float64) {
    x, y, z := col.Xyz()
    return XyzToLabWhiteRef(x, y, z, Xwref, Ywref, Zwref)
}

func Lab(l, a, b float64) Color {
    return Xyz(LabToXyz(l, a, b))
}

func LabWhiteRef(l, a, b, Xwref, Ywref, Zwref float64) Color {
    return Xyz(LabToXyzWhiteRef(l, a, b, Xwref, Ywref, Zwref))
}

func HappyColor() Color {
    H := rand.Float64() * 360.0
    S := 0.8+rand.Float64()*0.2
    V := 0.6+rand.Float64()*0.3

    return Hsv(H,S,V)
}
