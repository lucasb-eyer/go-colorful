package colorful

import "math"

// achromaticC is the Oklab chroma below which a color is treated as gray. Pure
// grays (and white/black) still carry a tiny non-zero b component from the
// matrix rounding, so an exact zero check is not enough to avoid dividing by a
// near-zero chroma when normalizing the hue.
const achromaticC = 1e-6

// Okhsl and Okhsv are two cylindrical color spaces built on top of Oklab by
// Björn Ottosson. Unlike Hsl/Hsv they are perceptually uniform, and unlike the
// plain OkLch cylinder their chroma is scaled so that saturation 1 always lands
// on the sRGB gamut boundary. This makes them well suited for color pickers.
//
// See https://bottosson.github.io/posts/colorpicker/ for the reference and a
// detailed description of the math below.
//
// The hue is returned in degrees in [0..360], while saturation and
// lightness/value are in [0..1].

// The Oklab conversions used here operate directly on linear sRGB, matching
// Ottosson's reference. This differs slightly from the XYZ-based OkLab helpers
// in colors.go, but keeps the gamut math (which is tuned for these exact
// coefficients) self-consistent.

func linearRgbToOklab(r, g, b float64) (L, a, bb float64) {
	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b

	l_ := math.Cbrt(l)
	m_ := math.Cbrt(m)
	s_ := math.Cbrt(s)

	L = 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	a = 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	bb = 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_
	return
}

func oklabToLinearRgb(L, a, b float64) (r, g, bb float64) {
	l_ := L + 0.3963377774*a + 0.2158037573*b
	m_ := L - 0.1055613458*a - 0.0638541728*b
	s_ := L - 0.0894841775*a - 1.2914855480*b

	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	r = +4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	g = -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	bb = -0.0041960863*l - 0.7034186147*m + 1.7076147010*s
	return
}

// toe applies the "toe" adjustment that makes Okhsl's lightness match the CIE L
// reference lightness. toeInv is its inverse.
func toe(x float64) float64 {
	const k1 = 0.206
	const k2 = 0.03
	const k3 = (1.0 + k1) / (1.0 + k2)
	return 0.5 * (k3*x - k1 + math.Sqrt((k3*x-k1)*(k3*x-k1)+4*k2*k3*x))
}

func toeInv(x float64) float64 {
	const k1 = 0.206
	const k2 = 0.03
	const k3 = (1.0 + k1) / (1.0 + k2)
	return (x*x + k1*x) / (k3 * (x + k2))
}

// computeMaxSaturation returns the maximum saturation S = C/L that still fits
// inside the sRGB gamut for a given hue. a and b must be normalized so that
// a*a + b*b == 1.
func computeMaxSaturation(a, b float64) float64 {
	var k0, k1, k2, k3, k4, wl, wm, ws float64

	switch {
	case -1.88170328*a-0.80936493*b > 1:
		// Red component
		k0, k1, k2, k3, k4 = 1.19086277, 1.76576728, 0.59662641, 0.75515197, 0.56771245
		wl, wm, ws = 4.0767416621, -3.3077115913, 0.2309699292
	case 1.81444104*a-1.19445276*b > 1:
		// Green component
		k0, k1, k2, k3, k4 = 0.73956515, -0.45954404, 0.08285427, 0.12541070, 0.14503204
		wl, wm, ws = -1.2684380046, 2.6097574011, -0.3413193965
	default:
		// Blue component
		k0, k1, k2, k3, k4 = 1.35733652, -0.00915799, -1.15130210, -0.50559606, 0.00692167
		wl, wm, ws = -0.0041960863, -0.7034186147, 1.7076147010
	}

	// Approximate max saturation using a polynomial.
	S := k0 + k1*a + k2*b + k3*a*a + k4*a*b

	// One step of Halley's method to refine.
	kl := 0.3963377774*a + 0.2158037573*b
	km := -0.1055613458*a - 0.0638541728*b
	ks := -0.0894841775*a - 1.2914855480*b

	l_ := 1.0 + S*kl
	m_ := 1.0 + S*km
	s_ := 1.0 + S*ks

	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	ldS := 3.0 * kl * l_ * l_
	mdS := 3.0 * km * m_ * m_
	sdS := 3.0 * ks * s_ * s_

	ldS2 := 6.0 * kl * kl * l_
	mdS2 := 6.0 * km * km * m_
	sdS2 := 6.0 * ks * ks * s_

	f := wl*l + wm*m + ws*s
	f1 := wl*ldS + wm*mdS + ws*sdS
	f2 := wl*ldS2 + wm*mdS2 + ws*sdS2

	return S - f*f1/(f1*f1-0.5*f*f2)
}

// lc holds a cusp: the most saturated, brightest point for a hue.
type lc struct {
	L, C float64
}

// st holds a cusp expressed as saturation/toe intercepts.
type st struct {
	S, T float64
}

// findCusp returns the cusp (highest chroma point in gamut) for a hue. a and b
// must be normalized so that a*a + b*b == 1.
func findCusp(a, b float64) lc {
	sCusp := computeMaxSaturation(a, b)

	r, g, bl := oklabToLinearRgb(1, sCusp*a, sCusp*b)
	lCusp := math.Cbrt(1.0 / math.Max(math.Max(r, g), bl))
	return lc{L: lCusp, C: lCusp * sCusp}
}

// findGamutIntersection finds t such that the point (L, C) = (L0*(1-t)+t*L1,
// t*C1) lies on the sRGB gamut boundary. a and b must be normalized.
func findGamutIntersection(a, b, L1, C1, L0 float64, cusp lc) float64 {
	var t float64
	if (L1-L0)*cusp.C-(cusp.L-L0)*C1 <= 0 {
		// Lower half.
		t = cusp.C * L0 / (C1*cusp.L + cusp.C*(L0-L1))
		return t
	}

	// Upper half: first intersect with the triangle.
	t = cusp.C * (L0 - 1.0) / (C1*(cusp.L-1.0) + cusp.C*(L0-L1))

	// Then one step of Halley's method to follow the curved boundary.
	dL := L1 - L0
	dC := C1

	kl := 0.3963377774*a + 0.2158037573*b
	km := -0.1055613458*a - 0.0638541728*b
	ks := -0.0894841775*a - 1.2914855480*b

	ldt := dL + dC*kl
	mdt := dL + dC*km
	sdt := dL + dC*ks

	L := L0*(1.0-t) + t*L1
	C := t * C1

	l_ := L + C*kl
	m_ := L + C*km
	s_ := L + C*ks

	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	ldt1 := 3 * ldt * l_ * l_
	mdt1 := 3 * mdt * m_ * m_
	sdt1 := 3 * sdt * s_ * s_

	ldt2 := 6 * ldt * ldt * l_
	mdt2 := 6 * mdt * mdt * m_
	sdt2 := 6 * sdt * sdt * s_

	rr := 4.0767416621*l - 3.3077115913*m + 0.2309699292*s - 1
	rr1 := 4.0767416621*ldt1 - 3.3077115913*mdt1 + 0.2309699292*sdt1
	rr2 := 4.0767416621*ldt2 - 3.3077115913*mdt2 + 0.2309699292*sdt2
	ur := rr1 / (rr1*rr1 - 0.5*rr*rr2)
	tr := -rr * ur

	gg := -1.2684380046*l + 2.6097574011*m - 0.3413193965*s - 1
	gg1 := -1.2684380046*ldt1 + 2.6097574011*mdt1 - 0.3413193965*sdt1
	gg2 := -1.2684380046*ldt2 + 2.6097574011*mdt2 - 0.3413193965*sdt2
	ug := gg1 / (gg1*gg1 - 0.5*gg*gg2)
	tg := -gg * ug

	bb := -0.0041960863*l - 0.7034186147*m + 1.7076147010*s - 1
	bb1 := -0.0041960863*ldt1 - 0.7034186147*mdt1 + 1.7076147010*sdt1
	bb2 := -0.0041960863*ldt2 - 0.7034186147*mdt2 + 1.7076147010*sdt2
	ub := bb1 / (bb1*bb1 - 0.5*bb*bb2)
	tb := -bb * ub

	if ur >= 0 {
		// keep tr
	} else {
		tr = math.MaxFloat64
	}
	if ug < 0 {
		tg = math.MaxFloat64
	}
	if ub < 0 {
		tb = math.MaxFloat64
	}

	t += math.Min(tr, math.Min(tg, tb))
	return t
}

func getSTMax(a, b float64, cusp lc) st {
	if cusp.L < 0 {
		cusp = findCusp(a, b)
	}
	return st{S: cusp.C / cusp.L, T: cusp.C / (1 - cusp.L)}
}

func getSTMid(a, b float64) st {
	s := 0.11516993 + 1.0/(7.44778970+4.15901240*b+
		a*(-2.19557347+1.75198401*b+
			a*(-2.13704948-10.02301043*b+
				a*(-4.24894561+5.38770819*b+4.69891013*a))))

	t := 0.11239642 + 1.0/(1.61320320-0.68124379*b+
		a*(0.40370612+0.90148123*b+
			a*(-0.27087943+0.61223990*b+
				a*(0.00299215-0.45399568*b-0.14661872*a))))

	return st{S: s, T: t}
}

// getCs returns the chroma control points (C0, Cmid, Cmax) used to shape the
// Okhsl saturation curve for a given lightness and hue.
func getCs(L, a, b float64) (c0, cmid, cmax float64) {
	cusp := findCusp(a, b)

	cmax = findGamutIntersection(a, b, L, 1, L, cusp)
	stMax := getSTMax(a, b, cusp)

	k := cmax / math.Min(L*stMax.S, (1-L)*stMax.T)

	stMid := getSTMid(a, b)
	ca := L * stMid.S
	cb := (1.0 - L) * stMid.T
	cmid = 0.9 * k * math.Sqrt(math.Sqrt(1.0/(1.0/(ca*ca*ca*ca)+1.0/(cb*cb*cb*cb))))

	ca = L * 0.4
	cb = (1.0 - L) * 0.8
	c0 = math.Sqrt(1.0 / (1.0/(ca*ca) + 1.0/(cb*cb)))
	return
}

// Okhsl returns the Okhsl representation of the color. The hue h is in degrees
// [0..360], and s and l are in [0..1].
func (col Color) Okhsl() (h, s, l float64) {
	lr, lg, lb := col.LinearRgb()
	L, a, b := linearRgbToOklab(lr, lg, lb)

	C := math.Sqrt(a*a + b*b)
	if C < achromaticC {
		return 0, 0, toe(L)
	}

	a_ := a / C
	b_ := b / C

	h = 0.5 + 0.5*math.Atan2(-b, -a)/math.Pi

	c0, cmid, cmax := getCs(L, a_, b_)

	const mid = 0.8
	const midInv = 1.25

	if C < cmid {
		k1 := mid * c0
		k2 := 1.0 - k1/cmid
		t := C / (k1 + k2*C)
		s = t * mid
	} else {
		k0 := cmid
		k1 := (1.0 - mid) * cmid * cmid * midInv * midInv / c0
		k2 := 1.0 - k1/(cmax-cmid)
		t := (C - k0) / (k1 + k2*(C-k0))
		s = mid + (1.0-mid)*t
	}

	l = toe(L)
	return h * 360, s, l
}

// Okhsl creates a Color from Okhsl coordinates. The hue h is in degrees, and s
// and l are in [0..1]. The result is clamped into the sRGB gamut.
func Okhsl(h, s, l float64) Color {
	if l >= 1 {
		return Color{1, 1, 1}
	} else if l <= 0 {
		return Color{0, 0, 0}
	}

	h /= 360
	a_ := math.Cos(2 * math.Pi * h)
	b_ := math.Sin(2 * math.Pi * h)
	L := toeInv(l)

	c0, cmid, cmax := getCs(L, a_, b_)

	const mid = 0.8
	const midInv = 1.25

	var C float64
	if s < mid {
		t := midInv * s
		k1 := mid * c0
		k2 := 1.0 - k1/cmid
		C = t * k1 / (1.0 - k2*t)
	} else {
		t := (s - mid) / (1 - mid)
		k0 := cmid
		k1 := (1.0 - mid) * cmid * cmid * midInv * midInv / c0
		k2 := 1.0 - k1/(cmax-cmid)
		C = k0 + t*k1/(1.0-k2*t)
	}

	r, g, b := oklabToLinearRgb(L, C*a_, C*b_)
	return LinearRgb(r, g, b).Clamped()
}

// Okhsv returns the Okhsv representation of the color. The hue h is in degrees
// [0..360], and s and v are in [0..1].
func (col Color) Okhsv() (h, s, v float64) {
	lr, lg, lb := col.LinearRgb()
	L, a, b := linearRgbToOklab(lr, lg, lb)

	C := math.Sqrt(a*a + b*b)
	if C < achromaticC {
		return 0, 0, toe(L)
	}

	a_ := a / C
	b_ := b / C

	h = 0.5 + 0.5*math.Atan2(-b, -a)/math.Pi

	stMax := getSTMax(a_, b_, lc{L: -1, C: -1})
	sMax := stMax.S
	tMax := stMax.T
	const s0 = 0.5
	k := 1 - s0/sMax

	t := tMax / (C + L*tMax)
	Lv := t * L
	Cv := t * C

	Lvt := toeInv(Lv)
	Cvt := Cv * Lvt / Lv

	rr, gg, bb := oklabToLinearRgb(Lvt, a_*Cvt, b_*Cvt)
	scaleL := math.Cbrt(1.0 / math.Max(math.Max(rr, gg), math.Max(bb, 0.0)))

	L = L / scaleL
	C = C / scaleL

	C = C * toe(L) / L
	L = toe(L)

	v = L / Lv
	s = (s0 + tMax) * Cv / (tMax*s0 + tMax*k*Cv)

	return h * 360, s, v
}

// Okhsv creates a Color from Okhsv coordinates. The hue h is in degrees, and s
// and v are in [0..1]. The result is clamped into the sRGB gamut.
func Okhsv(h, s, v float64) Color {
	h /= 360
	a_ := math.Cos(2 * math.Pi * h)
	b_ := math.Sin(2 * math.Pi * h)

	stMax := getSTMax(a_, b_, lc{L: -1, C: -1})
	sMax := stMax.S
	tMax := stMax.T
	const s0 = 0.5
	k := 1 - s0/sMax

	Lv := 1 - s*s0/(s0+tMax-tMax*k*s)
	Cv := s * tMax * s0 / (s0 + tMax - tMax*k*s)

	L := v * Lv
	C := v * Cv

	Lvt := toeInv(Lv)
	Cvt := Cv * Lvt / Lv

	Lnew := toeInv(L)
	if L != 0 {
		C = C * Lnew / L
	}
	L = Lnew

	rr, gg, bb := oklabToLinearRgb(Lvt, a_*Cvt, b_*Cvt)
	scaleL := math.Cbrt(1.0 / math.Max(math.Max(rr, gg), math.Max(bb, 0.0)))

	L = L * scaleL
	C = C * scaleL

	r, g, b := oklabToLinearRgb(L, C*a_, C*b_)
	return LinearRgb(r, g, b).Clamped()
}
