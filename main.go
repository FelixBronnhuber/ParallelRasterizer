package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	WIDTH  = 2560
	HEIGHT = 2560
)

var SCENE scene

type ray struct {
	origin, direction Vector
}

type sphere struct {
	origin Vector
	radius float64
	color  color
}

type scene struct {
	spheres []sphere
}

type color struct {
	r, g, b uint8
}

func main() {
	// Sphere intersection test:
	/*
		SCENE = scene{makeSphereArray(
			sphere{Vector{100, 70, 0}, 200, color{255, 0, 0},},
			sphere{Vector{70, -70, 0}, 200, color{0, 0, 255},},
			sphere{Vector{70, 70, -70}, 200, color{0, 255, 0},},
		)}
	*/

	SCENE.randomSphereTest(5000)
	saveImage(renderImage())
}

// Saves the image as a ".png" file to the out folder. Adds a time stamp to the file name.
func saveImage(_image image.Image) {
	name := "out/render_" + time.Now().Format("2006_01_02_15_04_05") + ".png"
	outputFile, err := os.Create(name)
	err = png.Encode(outputFile, _image)
	err = outputFile.Close()
	if err != nil {
		fmt.Println("Error Occurred")
		fmt.Println(err)
	}
}

// Renders the image from the scene.
func renderImage() image.Image {
	// To measure the time it took to render the image.
	start := time.Now()
	var _image = image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	// For a correct field of view:
	d := (float64(WIDTH / 2)) / (math.Tan(math.Pi / 4))
	// Raycast from origin (at [d, 0, 0]) through the particular pixel on the view port.
	rays := make([]ray, WIDTH*HEIGHT)
	c := 0
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			c = y*HEIGHT + x
			rays[c] = ray{
				origin:    Vector{d, 0, 0},
				direction: Vector{-d, float64(x - WIDTH/2), float64(y - HEIGHT/2)},
			}
		}
	}

	// Hopefully concurrent raycasts.
	prepRays := rayPipeline(rays)
	getC := raycastWorker(prepRays)

	// To stare at while it's running.
	// go func() {
	// 	for {
	// 		fmt.Print("Working ")
	// 		for i := 0; i < 20; i++ {
	// 			fmt.Print(".")
	// 			time.Sleep(time.Millisecond * 2000)
	// 		}
	// 		fmt.Println()
	// 	}
	// }()

	// Debug output: Displays the amount of active goroutines and log. CPU cores.
	go func() {
		for {
			fmt.Println()
			fmt.Println(time.Now().Format(time.StampMilli))
			fmt.Println(runtime.NumGoroutine(), "Gophers are working.")
			fmt.Println(runtime.NumCPU(), "logical CPU Cores are available.")
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	// Converts the output color channel to an array.
	pixels := make([]color, WIDTH*HEIGHT)

	i := 0
	for n := range getC {
		pixels[i] = n
		i++
	}

	// Converts the array to an image file.
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			// TODO: Not working properly for non square, rectangular images.
			pO := _image.PixOffset(x, y)
			pI := x*WIDTH + y
			_image.Pix[pO] = pixels[pI].r
			_image.Pix[pO+1] = pixels[pI].g
			_image.Pix[pO+2] = pixels[pI].b
			_image.Pix[pO+3] = 255
		}
	}

	fmt.Println("RENDER TIME:", time.Since(start))
	return _image
}

// Populates a chanel with the rays...
func rayPipeline(rays []ray) <-chan ray {
	var out = make(chan ray)
	go func() {
		for n := range rays {
			out <- rays[n]
		}
		close(out)
	}()
	return out
}

// ... passes it on to the Worker func.
func raycastWorker(in <-chan ray) <-chan color {
	var wg sync.WaitGroup
	wg.Add(12)
	var out = make(chan color)
	go func() {
		for i := 0; i < 12; i++ {
			go func() {
				defer wg.Done()
				for n := range in {
					out <- n.intersectsWith().color
				}
			}()
		}
		wg.Wait()
		close(out)
	}()
	return out
}

// Solves the equation for Line–Sphere intersection. Returns the closest Sphere it intersects with.
func (r ray) intersectsWith() sphere {
	abs, closestSphere := math.MaxFloat64, sphere{}

	for n := range SCENE.spheres {
		s := SCENE.spheres[n]
		v := Vector{
			X1: r.origin.X1 - s.origin.X1,
			X2: r.origin.X2 - s.origin.X2,
			X3: r.origin.X3 - s.origin.X3,
		}
		disc := math.Pow(2*r.direction.GetDotProduct(v), 2) -
			(4*r.direction.GetVecSquared())*(v.GetVecSquared()-
				math.Pow(s.radius, 2))
		// There is no solution:
		if disc < 0 {
			n++
		}
		div := 2 * r.direction.GetVecSquared()
		b := -2 * r.direction.GetDotProduct(v)
		lambda0 := (b + math.Sqrt(disc)) / div
		lambda1 := (b - math.Sqrt(disc)) / div
		p0 := r.getPointForLambda(lambda0)
		p1 := r.getPointForLambda(lambda1)
		abs0 := p0.DifferenceVector(r.origin).GetAbs()
		abs1 := p1.DifferenceVector(r.origin).GetAbs()
		if abs0 < abs1 {
			if abs0 < abs {
				abs = abs0
				closestSphere = s
			}
		} else {
			if abs1 < abs {
				abs = abs1
				closestSphere = s
			}
		}
	}

	return closestSphere
}

// Returns the Parameters of the function as an array of spheres.
func makeSphereArray(spheres ...sphere) []sphere {
	return spheres
}

// Returns a Point on the line for a given value.
func (r ray) getPointForLambda(l float64) Vector {
	v := Vector{
		X1: r.origin.X1 + (l * r.direction.X1),
		X2: r.origin.X2 + (l * r.direction.X2),
		X3: r.origin.X3 + (l * r.direction.X3),
	}
	return v
}

// Populates the Scene with Spheres of random color at random positions.
func (s *scene) randomSphereTest(n int) {
	seed := time.Now().Unix()
	rand.Seed(seed)
	fmt.Println("Seed: ", seed)

	s.spheres = make([]sphere, n)
	for i := 0; i < n; i++ {
		s.spheres[i] = sphere{
			origin: Vector{
				rand.Float64()*WIDTH - (WIDTH / 2),
				rand.Float64()*WIDTH - (WIDTH / 2),
				rand.Float64()*WIDTH - (WIDTH / 2),
			},
			radius: 10,
			color: color{
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
				uint8(rand.Intn(255)),
			},
		}
	}
}
