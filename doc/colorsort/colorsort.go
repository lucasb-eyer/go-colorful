// This program generates an example of go-colorful's Sorted function.  It
// produces an image with three stripes of color.  The first is unsorted.
// The second is sorted in the CIE-L\*C\*h° space, ordered primarily by
// lightness, then by hue angle, and finally by chroma.  The third is
// sorted using colorful.Sorted.

package main

import (
	"image"
	"image/png"
	"math/rand"
	"os"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
)

// randomColors produces a slice of random colors.
func randomColors(n int) []colorful.Color {
	cs := make([]colorful.Color, n)
	for i := range cs {
		cs[i] = colorful.Color{
			R: rand.Float64(),
			G: rand.Float64(),
			B: rand.Float64(),
		}
	}
	return cs
}

// drawStripes creates an image with three sets of stripes.
func drawStripes(cs1, cs2, cs3 []colorful.Color, ht, sep int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, len(cs1), 3*ht+2*sep))
	for c := range cs1 {
		for r := 0; r < ht; r++ {
			img.Set(c, r, cs1[c].Clamped())
			img.Set(c, r+ht+sep, cs2[c].Clamped())
			img.Set(c, r+(ht+sep)*2, cs3[c].Clamped())
		}
	}
	return img
}

// writeImage writes an image to disk.
func writeImage(fn string, img image.Image) {
	w, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png.Encode(w, img)
	if err != nil {
		panic(err)
	}
}

func main() {
	n := 512
	rand.Seed(8675309)
	cs1 := randomColors(n)
	cs2 := make([]colorful.Color, n)
	copy(cs2, cs1)
	sort.Slice(cs2, func(i, j int) bool {
		l1, c1, h1 := cs2[i].LuvLCh()
		l2, c2, h2 := cs2[j].LuvLCh()
		if l1 != l2 {
			return l1 < l2
		}
		if h1 != h2 {
			return h1 < h2
		}

		if c1 != c2 {
			return c1 < c2
		}
		return false
	})
	cs3 := colorful.Sorted(cs1)
	img := drawStripes(cs1, cs2, cs3, 64, 16)
	writeImage("colorsort.png", img)
}
