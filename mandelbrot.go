package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"time"
	"sync"
	"github.com/lucasb-eyer/go-colorful"
)

var smoothingConstant = 1 / math.Log2(math.Log2(cmplx.Abs(6)))

func main() {
	start := time.Now()
	centerX, centerY, width, height, zoomLevel := -0.745428, 0.113009, 16.0, 9.0, 2e5
	resolution, maxIterations := 1080, 1000
	bounds := image.Rect(0, 0, int(float64(resolution)*width/height), resolution)
	img := image.NewRGBA(bounds)
	realSize, imagSize := width/zoomLevel, height/zoomLevel
	reals, imags := make([]float64, bounds.Max.X, bounds.Max.X), make([]float64, bounds.Max.Y, bounds.Max.Y)
	for x := 0; x < bounds.Max.X; x++ {
		reals[x] = centerX - realSize/2 + realSize*(float64(x)+0.5)/float64(bounds.Max.X)
	}
	for y := 0; y < bounds.Max.Y; y++ {
		imags[y] = centerY + imagSize/2 - imagSize*(float64(y)+0.5)/float64(bounds.Max.Y)
	}
	var wg sync.WaitGroup
	for x, re := range reals {
		wg.Add(1)
		go func(x int, re float64) {
			for y, im := range imags {
				z, c, i := complex(0, 0), complex(re, im), 0
				for cmplx.Abs(z) < 2 && i < maxIterations {
					z = z*z + c
					i++
				}
				if i < maxIterations {
					distance := (float64(i+1) - math.Log2(math.Log2(cmplx.Abs(z)))*smoothingConstant) / float64(maxIterations)
					r, g, b := colorful.Hcl(distance*360, distance, distance).RGB255()
					img.Set(x, y, color.RGBA{r, g, b, 255})
				} else {
					img.Set(x, y, color.RGBA{0, 0, 0, 255})
				}
			}
			wg.Done()
		}(x, re)
	}
	wg.Wait()
	file, _ := os.Create("mandelbrot.png")
	png.Encode(file, img)
	file.Close()
	fmt.Println(time.Now().Sub(start))
}
