go-colorful
===========
A library for playing with colors in go (golang).

Why?
====
I love games. I make games. I love detail and I get lost in detail.
One such detail popped up during the development of [Memory Which Does Not Suck](https://github.com/lucasb-eyer/mwdns/),
when we wanted the server to assign the players random colors. Sometimes
two players got very similar colors, which bugged me. The very same evening,
[I want hue](http://tools.medialab.sciences-po.fr/iwanthue/) was the top post
on HackerNews' frontpage and showed me how to Do It Right (tm). Last but not
least, there was no library for handling color spaces available in go.

What?
=====
Go-Colorful stores colors in RGB and provides methods from converting these to various color-spaces. Currently supported colorspaces are:

- **RGB:** All three of Red, Green and Blue in [0..1].
- **HSV:** Hue in [0..360], Saturation and Value in [0..1]. You probably shouldn't use this.
- **Hex RGB:** The "internet" color format, as in #FF00FF.
- **Linear RGB:** See [gamma correct rendering](http://www.sjbrown.co.uk/2004/05/14/gamma-correct-rendering/).
- **CIE XYZ:** CIE's standard color space, almost in [0..1].
- **L\*a\*b\*:** A *perceptually uniform* color space, i.e. distances are meaningful. L\* in [0..1] and a\*, b\* almost in [-1..1].

For the colorspaces where it makes sense (XYZ, Lab), the
[D65](http://en.wikipedia.org/wiki/Illuminant_D65) is used as reference white
by default but methods for using your own reference white are provided.

Unit-tests are provided.

What not (yet)?
===============
There are a few features which are currently missing and might be useful.
I just haven't implemented them yet because I didn't have the need for it.
Pull requests welcome.

- Functions for palette generation. (This is in the making.)
- Functions for computing distances in various spaces. Note that this seems to be [a whole science on its own](http://mmir.doc.ic.ac.uk/mmir2005/CameraReadyMissaoui.pdf).
- Functions for interpolation.
- The [HCL](http://vis4.net/blog/posts/avoid-equidistant-hsv-colors/) color space.
- The [CIELUV](http://en.wikipedia.org/wiki/CIELUV) color space.

So which colorspace should I use?
=================================
It depends on what you want to do. I think the folks from *I want hue* are
on-spot when they say that RGB fits to how *screens produce* color, CIE L\*a\*b\*
fits how *humans perceive* color and HCL fits how *humans think* colors.

How?
====

### Installing
Installing the library is as easy as

```bash
$ go get github.com/lucasb-eyer/go-colorful
```

The package can then be used through an

```go
import "github.com/lucasb-eyer/go-colorful"
```

### Basic usage

Create a beautiful blue color using different source space:

```go
c := colorful.Color{0.313725, 0.478431, 0.721569}
c := colorful.Hex("#517AB8")
c := colorful.Hsv(216.0, 0.56, 0.722)
c := colorful.Xyz(0.189165, 0.190837, 0.480248)
c := colorful.Lab(0.507850, 0.040585,-0.370945)
Printf("RGB values: %v, %v, %v", c.R, c.G, c.B)
```

And then converting this color back into various color spaces:

```go
hex := c.Hex()
h, s, v := c.Hsv()
x, y, z := c.Xyz()
l, a, b := c.Lab()
```

### Comparing colors
In the RGB color space, the Euclidian distance between colors *doesn't* correspond
to visual/perceptual distance. This means that two pairs of colors which have the
same distance in RGB space can look much further apart. This is fixed by the
L\*a\*b\* color space. Thus you should only compare colors in that space.

![Color distance comparison](doc/colordist.png)

The two colors shown on the top look much more different than the two shown on
the bottom. Still, in RGB space, their distance is the same.
Here is a little example program which shows the distances between the top two
and bottom two colors in RGB and in L\*a\*b\* space. You can find it in `doc/colordist.go`.

```go
package main

import "github.com/lucasb-eyer/go-colorful"

func main() {
    c1a := colorful.Color{150.0/255.0, 10.0/255.0, 150.0/255.0}
    c1b := colorful.Color{ 53.0/255.0, 10.0/255.0, 150.0/255.0}
    c2a := colorful.Color{10.0/255.0, 150.0/255.0, 50.0/255.0}
    c2b := colorful.Color{99.9/255.0, 150.0/255.0, 10.0/255.0}

    fmt.Printf("DistanceRgb: c1: %v and c2: %v\n", c1a.DistanceRgb(c1b), c2a.DistanceRgb(c2b))
    fmt.Printf("DistanceLab: c1: %v and c2: %v\n", c1a.DistanceLab(c1b), c2a.DistanceLab(c2b))
}
```

Running the above program shows that you should always prefer L\*a\*b\* distance:

```bash
$ go run colordist.go
DistanceRgb: 0.3803921568627451 and 0.3858713931171159
DistanceLab: 0.32043169229603413 and 0.24397340448513863
```

Note that `AlmostEqualRgb` is provided mainly for (unit-)testing purposes. Use
it only if you really know what you're doing. It will eat your cat.

### Getting random palettes/colors
TODO

### Using linear RGB for computations
There are two methods for transforming RGB<->Linear RGB: a fast and almost precise one,
and a slow and precise one.

```go
r, g, b := colorful.Hex("#FF0000").FastLinearRgb()
```

TODO: describe some more.

### Want to use some other reference point?

```go
c := colorful.LabWhiteRef(0.507850, 0.040585,-0.370945, colorful.D50)
l, a, b := c.LabWhiteRef(colorful.D50)
```

FAQ
===

### Q: I get all f!@#ed up values! Your library sucks!
A: You probably provided values in the wrong range. For example, RGB values are
expected to reside between 0 and 1, *not* between 0 and 255. Normalize your colors.

License: MIT
============
Copyright (c) 2013 Lucas Beyer

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

