package main

import "fmt"
import "github.com/lucasb-eyer/go-colorful"
import "image"
import "image/draw"
// import "image/color"
import "image/png"
import "math"
import "os"

type GradientTable []struct {
    Col colorful.Color
    Pos float64
}

func (self GradientTable) GetInerpolatedColorFor(t float64) colorful.Color {
    for i := 0 ; i < len(self) - 1 ; i++ {
        c1 := self[i]
        c2 := self[i+1]
        if c1.Pos <= t && t <= c2.Pos {
            t := (t - c1.Pos)/(c2.Pos - c1.Pos)
            return c1.Col.BlendHcl(c2.Col, t).Clamped()
        }
    }

    return self[len(self)-1].Col
}

func (self GradientTable) GetColorFor(t float64) colorful.Color {
    for i := 0 ; i < len(self) - 1 ; i++ {
        // Note: This relies on the fact that they are sorted from smallest to largest.
        if t < (self[i].Pos + self[i+1].Pos)*0.5 {
            return self[i].Col
        }
    }

    return self[len(self)-1].Col
}

func savepng(img image.Image, fname string) {
    toimg, err := os.Create(fname)
    if err != nil {
        fmt.Printf("Error: %v", err)
        return
    }
    defer toimg.Close()

    png.Encode(toimg, img)
}

func MustParseHex(s string) colorful.Color {
    c, err := colorful.Hex(s)
    if err != nil {
        panic("MustParseHex: " + err.Error())
    }
    return c
}

func main() {
    // The "keypoints" of the gradient.
    colorbrewer_spectral := GradientTable{
        {MustParseHex("#9e0142"), 0.0},
        {MustParseHex("#d53e4f"), 0.1},
        {MustParseHex("#f46d43"), 0.2},
        {MustParseHex("#fdae61"), 0.3},
        {MustParseHex("#fee090"), 0.4},
        {MustParseHex("#ffffbf"), 0.5},
        {MustParseHex("#e6f598"), 0.6},
        {MustParseHex("#abdda4"), 0.7},
        {MustParseHex("#66c2a5"), 0.8},
        {MustParseHex("#3288bd"), 0.9},
        {MustParseHex("#5e4fa2"), 1.0},
    }

    h := 1024
    w := 40
    img_linear := image.NewRGBA(image.Rect(0,0,w,h))
    img_square := image.NewRGBA(image.Rect(0,0,w,h))
    img_sqroot := image.NewRGBA(image.Rect(0,0,w,h))

    fmt.Println("const unsigned char _g_heatmap_colorscheme_colorbrewer_spectral_mixed_exp_data[] = {")
    fmt.Println("0, 0, 0, 0,")
    for y := h-1 ; y >= 0 ; y-- {
        t := float64(y)/float64(h)
        c1 := colorbrewer_spectral.GetInerpolatedColorFor(t)
        c2 := colorbrewer_spectral.GetColorFor(t)
        draw.Draw(img_linear, image.Rect(0, y, w, y+1), &image.Uniform{c1.BlendRgb(c2, 0.2)}, image.ZP, draw.Src)
        c1 = colorbrewer_spectral.GetInerpolatedColorFor(t*t*t*t*t*t*t*t*t*t)
        c2 = colorbrewer_spectral.GetColorFor(t*t*t*t*t*t*t*t*t*t)
        draw.Draw(img_square, image.Rect(0, y, w, y+1), &image.Uniform{c1.BlendRgb(c2, 0.2)}, image.ZP, draw.Src)
        r, g, b := c1.BlendRgb(c2, 0.2).RGB255()
        fmt.Printf("%v, %v, %v, 255,\n", r, g, b)
        c1 = colorbrewer_spectral.GetInerpolatedColorFor(math.Sqrt(t))
        c2 = colorbrewer_spectral.GetColorFor(math.Sqrt(t))
        draw.Draw(img_sqroot, image.Rect(0, y, w, y+1), &image.Uniform{c1.BlendRgb(c2, 0.2)}, image.ZP, draw.Src)
    }
    fmt.Println("};")

    savepng(img_linear, "gradient_linear.png")
    savepng(img_square, "gradient_square.png")
    savepng(img_sqroot, "gradient_sqroot.png")
}
