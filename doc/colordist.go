package main

import "fmt"
//import "github.com/lucasb-eyer/go-colorful"
import "go-colorful"

func main() {
    c1a := colorful.Color{150.0/255.0, 10.0/255.0, 150.0/255.0}
    c1b := colorful.Color{ 53.0/255.0, 10.0/255.0, 150.0/255.0}
    c2a := colorful.Color{10.0/255.0, 150.0/255.0, 50.0/255.0}
    c2b := colorful.Color{99.9/255.0, 150.0/255.0, 10.0/255.0}

    fmt.Printf("DistanceRgb: %v and %v\n", c1a.DistanceRgb(c1b), c2a.DistanceRgb(c2b))
    fmt.Printf("DistanceLab: %v and %v\n", c1a.DistanceLab(c1b), c2a.DistanceLab(c2b))
}
