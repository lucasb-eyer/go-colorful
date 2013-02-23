// The colorful package provides all kinds of functions for working with colors.
package colorful

import(
    "fmt"
    "math"
)

// A color is stored internally using sRGB (standard RGB) values in the range 0-1
type Color struct {
    R, G, B float64
}

// Implement the Go color.Color interface.
func (col Color) RGBA() (r, g, b, a uint32) {
    r = uint32(col.R*65535.0)
    g = uint32(col.G*65535.0)
    b = uint32(col.B*65535.0)
    a = 0xFFFF
    return
}

// This is the tolerance used when comparing colors using AlmostEqualRgb.
const Delta = 1.0/255.0

// This is the default reference white point.
var D65 = [3]float64{0.95043, 1.00000, 1.08890}

// And another one.
var D50 = [3]float64{0.96421, 1.00000, 0.82519}

// Checks whether the color exists in RGB space, i.e. all values are in [0..1]
func (c Color) IsValid() bool {
    return 0.0 <= c.R && c.R <= 1.0 &&
           0.0 <= c.G && c.G <= 1.0 &&
           0.0 <= c.B && c.B <= 1.0
}

func sq(v float64) float64 {
    return v * v;
}

func cub(v float64) float64 {
    return v * v * v;
}

// DistanceRgb computes the distance between two colors in RGB space.
// This is not a good measure! Rather do it in Lab space.
func (c1 Color) DistanceRgb(c2 Color) float64 {
    return math.Sqrt(sq(c1.R-c2.R) + sq(c1.G-c2.G) + sq(c1.B-c2.B))
}

// Check for equality between colors within the tolerance Delta (1/255).
func (c1 Color) AlmostEqualRgb(c2 Color) bool {
    return math.Abs(c1.R - c2.R) +
           math.Abs(c1.G - c2.G) +
           math.Abs(c1.B - c2.B) < 3.0*Delta
}

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
    return
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
func Hex(scol string) (Color, error) {
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
// http://www.brucelindbloom.com/Eqn_RGB_to_XYZ.html

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
    return
}

// FastLinearRgb is much faster than and almost as accurate as LinearRgb.
func (col Color) FastLinearRgb() (r, g, b float64) {
    r = math.Pow(col.R, 2.2)
    g = math.Pow(col.G, 2.2)
    b = math.Pow(col.B, 2.2)
    return
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
    r =  3.2404542*x - 1.5371385*y - 0.4985314*z
    g = -0.9692660*x + 1.8760108*y + 0.0415560*z
    b =  0.0556434*x - 0.2040259*y + 1.0572252*z
    return
}

func LinearRgbToXyz(r, g, b float64) (x, y, z float64) {
    x = 0.4124564*r + 0.3575761*g + 0.1804375*b
    y = 0.2126729*r + 0.7151522*g + 0.0721750*b
    z = 0.0193339*r + 0.1191920*g + 0.9503041*b
    return
}

/// XYZ ///
///////////
// http://www.sjbrown.co.uk/2004/05/14/gamma-correct-rendering/

func (col Color) Xyz() (x, y, z float64) {
    return LinearRgbToXyz(col.LinearRgb())
}

func Xyz(x, y, z float64) Color {
    return LinearRgb(XyzToLinearRgb(x, y, z))
}

/// L*a*b* ///
//////////////
// http://en.wikipedia.org/wiki/Lab_color_space#CIELAB-CIEXYZ_conversions
// For L*a*b*, we need to L*a*b*<->XYZ->RGB and the first one is device dependent.

func lab_f(t float64) float64 {
    if t > 6.0/29.0 * 6.0/29.0 * 6.0/29.0 {
        return math.Cbrt(t)
    }
    return t/3.0 * 29.0/6.0 * 29.0/6.0 + 4.0/29.0
}

func XyzToLab(x, y, z float64) (l, a, b float64) {
    // Use D65 white as reference point by default.
    // http://www.fredmiranda.com/forum/topic/1035332
    // http://en.wikipedia.org/wiki/Standard_illuminant
    return XyzToLabWhiteRef(x, y, z, D65)
}

func XyzToLabWhiteRef(x, y, z float64, wref [3]float64) (l, a, b float64) {
    fy := lab_f(y/wref[1])
    l = 1.16*fy - 0.16
    a = 5.0*(lab_f(x/wref[0]) - fy)
    b = 2.0*(fy - lab_f(z/wref[2]))
    return
}

func lab_finv(t float64) float64 {
    if t > 6.0/29.0 {
        return t * t * t
    }
    return 3.0 * 6.0/29.0 * 6.0/29.0 * (t - 4.0/29.0)
}

func LabToXyz(l, a, b float64) (x, y, z float64) {
    // D65 white (see above).
    return LabToXyzWhiteRef(l, a, b, D65)
}

func LabToXyzWhiteRef(l, a, b float64, wref [3]float64) (x, y, z float64) {
    l2 := (l + 0.16) / 1.16
    x = wref[0] * lab_finv(l2 + a/5.0)
    y = wref[1] * lab_finv(l2)
    z = wref[2] * lab_finv(l2 - b/2.0)
    return
}

// Converts the given color to CIE L*a*b* space using D65 as reference white.
func (col Color) Lab() (l, a, b float64) {
    return XyzToLab(col.Xyz())
}

// Converts the given color to CIE L*a*b* space, taking into account
// a given reference white. (i.e. the monitor's white)
func (col Color) LabWhiteRef(wref [3]float64) (l, a, b float64) {
    x, y, z := col.Xyz()
    return XyzToLabWhiteRef(x, y, z, wref)
}

// Generates a color by using data given in CIE L*a*b* space using D65 as reference white.
func Lab(l, a, b float64) Color {
    return Xyz(LabToXyz(l, a, b))
}

// Generates a color by using data given in CIE L*a*b* space, taking
// into account a given reference white. (i.e. the monitor's white)
func LabWhiteRef(l, a, b float64, wref [3]float64) Color {
    return Xyz(LabToXyzWhiteRef(l, a, b, wref))
}

// DistanceLab is a good measure of visual similarity between two colors!
// A result of 0 would mean identical colors, while a result of 1 or higher
// means the colors differ a lot.
func (c1 Color) DistanceLab(c2 Color) float64 {
    l1, a1, b1 := c1.Lab()
    l2, a2, b2 := c2.Lab()
    return math.Sqrt(sq(l1-l2) + sq(a1-a2) + sq(b1-b2))
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

/// L*u*v* ///
//////////////
// http://en.wikipedia.org/wiki/CIELUV#XYZ_.E2.86.92_CIELUV_and_CIELUV_.E2.86.92_XYZ_conversions
// For L*u*v*, we need to L*u*v*<->XYZ<->RGB and the first one is device dependent.

func XyzToLuv(x, y, z float64) (l, a, b float64) {
    // Use D65 white as reference point by default.
    // http://www.fredmiranda.com/forum/topic/1035332
    // http://en.wikipedia.org/wiki/Standard_illuminant
    return XyzToLuvWhiteRef(x, y, z, D65)
}

func XyzToLuvWhiteRef(x, y, z float64, wref [3]float64) (l, u, v float64) {
    if y/wref[1] <= 6.0/29.0 * 6.0/29.0 * 6.0/29.0 {
        l = y/wref[1] * 29.0/3.0 * 29.0/3.0 * 29.0/3.0
    } else {
        l = 1.16 * math.Cbrt(y/wref[1]) - 0.16
    }
    ubis, vbis := xyz_to_uv(x, y, z)
    un, vn := xyz_to_uv(wref[0], wref[1], wref[2])
    u = 13.0*l * (ubis - un)
    v = 13.0*l * (vbis - vn)
    return
}

// For this part, we do as R's graphics.hcl does, not as wikipedia does.
// Or is it the same?
func xyz_to_uv(x, y, z float64) (u, v float64) {
    denom := x + 15.0*y + 3.0*z
    if denom == 0.0 {
        u, v = 0.0, 0.0
    } else {
        u = 4.0*x/denom
        v = 9.0*y/denom
    }
    return
}

func LuvToXyz(l, u, v float64) (x, y, z float64) {
    // D65 white (see above).
    return LuvToXyzWhiteRef(l, u, v, D65)
}

func LuvToXyzWhiteRef(l, u, v float64, wref [3]float64) (x, y, z float64) {
    //y = wref[1] * lab_finv((l + 0.16) / 1.16)
    if l <= 0.08 {
        y = wref[1] * l * 100.0 * 3.0/29.0 * 3.0/29.0 * 3.0/29.0
    } else {
        y = wref[1] * cub((l+0.16)/1.16)
    }
    un, vn := xyz_to_uv(wref[0], wref[1], wref[2])
    if l != 0.0 {
        ubis := u/(13.0*l) + un
        vbis := v/(13.0*l) + vn
        x = y*9.0*ubis/(4.0*vbis)
        z = y*(12.0-3.0*ubis-20.0*vbis)/(4.0*vbis)
    } else {
        x, y = 0.0, 0.0
    }
    return
}

// Converts the given color to CIE L*u*v* space using D65 as reference white.
// L* is in [0..1] and both u* and v* are in about [-1..1]
func (col Color) Luv() (l, u, v float64) {
    return XyzToLuv(col.Xyz())
}

// Converts the given color to CIE L*u*v* space, taking into account
// a given reference white. (i.e. the monitor's white)
// L* is in [0..1] and both u* and v* are in about [-1..1]
func (col Color) LuvWhiteRef(wref [3]float64) (l, u, v float64) {
    x, y, z := col.Xyz()
    return XyzToLuvWhiteRef(x, y, z, wref)
}

// Generates a color by using data given in CIE L*u*v* space using D65 as reference white.
// L* is in [0..1] and both u* and v* are in about [-1..1]
func Luv(l, u, v float64) Color {
    return Xyz(LuvToXyz(l, u, v))
}

// Generates a color by using data given in CIE L*u*v* space, taking
// into account a given reference white. (i.e. the monitor's white)
// L* is in [0..1] and both u* and v* are in about [-1..1]
func LuvWhiteRef(l, u, v float64, wref [3]float64) Color {
    return Xyz(LuvToXyzWhiteRef(l, u, v, wref))
}

/// HCL ///
///////////
// HCL is nothing else than L*a*b* in cylindrical coordinates!
// (this was wrong on English wikipedia, I fixed it, let's hope the fix stays.)
// But it is widely popular since it is a "correct HSV"
// http://www.hunterlab.com/appnotes/an09_96a.pdf

// Converts the given color to HCL space using D65 as reference white.
// H values are in [0..360], C and L values are in [0..1] although C can overshoot 1.0
func (col Color) Hcl() (h, c, l float64) {
    return col.HclWhiteRef(D65)
}

// Converts the given color to HCL space, taking into account
// a given reference white. (i.e. the monitor's white)
// H values are in [0..360], C and L values are in [0..1]
func (col Color) HclWhiteRef(wref [3]float64) (h, c, l float64) {
    L, a, b := col.LabWhiteRef(wref)

    // Oops, floating point workaround necessary if a ~= b and both are very small (i.e. almost zero).
    if math.Abs(b - a) > 1e-4 && b > 1e-4 {
        h = math.Mod(57.29577951308232087721*math.Atan2(b, a) + 360.0, 360.0) // Rad2Deg
    } else {
        h = 0.0
    }
    c = math.Sqrt(sq(a) + sq(b))
    l = L
    return
}

// Generates a color by using data given in HCL space using D65 as reference white.
// H values are in [0..360], C and L values are in [0..1]
func Hcl(h, c, l float64) Color {
    return HclWhiteRef(h, c, l, D65)
}

// Generates a color by using data given in HCL space, taking
// into account a given reference white. (i.e. the monitor's white)
// H values are in [0..360], C and L values are in [0..1]
func HclWhiteRef(h, c, l float64, wref [3]float64) Color {
    H := 0.01745329251994329576*h // Deg2Rad
    a := c*math.Cos(H)
    b := c*math.Sin(H)
    return LabWhiteRef(l, a, b, wref)
}

