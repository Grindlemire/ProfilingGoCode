package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime/trace"
	"strings"
	"sync"

	log "github.com/cihub/seelog"
	"github.com/grindlemire/seezlog"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/profile"
)

// Opts ...
type Opts struct {
	File         string  `short:"f" long:"file" default:"output.png" description:"File name to output to"`
	Complexity   int     `short:"c" long:"complexity" default:"4" description:"Complexity of the fractal"`
	MaxIteration int     `short:"i" long:"maxIterations" default:"1000" description:"Max number of iterations to run in fractal calculation"`
	MoveX        float64 `short:"x" default:"0" description:"x movement"`
	MoveY        float64 `short:"y" default:"0" description:"y movement"`
	Zoom         float64 `short:"z" default:"1" description:"zoom level"`
	Width        int     `long:"width" default:"2048" description:"width of image"`
	Height       int     `long:"height" default:"2048" description:"height of image"`
	Mem          bool    `long:"mem" description:"memory profile"`
	CPU          bool    `long:"cpu" description:"cpu profile"`
	Trace        bool    `long:"trace" description:"trace profile"`
	Block        bool    `long:"block" description:"block profile"`
}

var opts Opts
var parser = flags.NewParser(&opts, flags.Default)

func main() {
	logger, err := seezlog.SetupConsoleLogger(seezlog.Info)
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		exit(1)
	}
	err = log.ReplaceLogger(logger)
	if err != nil {
		fmt.Printf("Error replacing logger: %s\n", err)
		exit(1)
	}
	defer log.Flush()

	_, err = parser.Parse()
	if err != nil {
		if !isUsage(err) {
			log.Error("Error parsing arguments: ", err)
			exit(1)
		} else {
			exit(0)
		}
	}

	// LaunchServer()

	f, err := os.Create(opts.File)
	if err != nil {
		log.Error("Error creating file: ", err)
		exit(1)
	}

	p := getProfiler()
	if p != nil {
		defer profile.Start(p, profile.ProfilePath("./")).Stop()
	}

	if opts.Trace {
		f, err := os.Create("./out.trace")
		if err != nil {
			log.Error("Error creating file: ", err)
			exit(1)
		}
		trace.Start(f)
		defer trace.Stop()
	}

	// img, err := executeAlgorithm(opts)
	// img, err := executeColumnParallelAlgorithm(opts)
	img, err := executeBufferedColumnWorkersAlgorithm(opts)
	// img, err := executeWorkersAlgorithm(opts)
	// img, err := executeBufferedWorkersAlgorithm(opts)
	// img, err := executePixelParallelAlgorithm(opts)
	if err != nil {
		log.Error("Error executing algorithm: ", err)
		exit(1)
	}

	err = png.Encode(f, img)
	if err != nil {
		log.Error("Error encoding fractal to file: ", err)
		exit(1)
	}

}

func getProfiler() func(p *profile.Profile) {
	if opts.CPU {
		return profile.CPUProfile
	}

	if opts.Mem {
		return profile.MemProfile
	}

	if opts.Block {
		return profile.BlockProfile
	}

	return nil
}

func executeAlgorithm(opts Opts) (img image.Image, err error) {
	m := createPNG(opts)
	for i := 0; i < opts.Width; i++ {
		for j := 0; j < opts.Height; j++ {
			m.Set(i, j, getMandelbrotColor(i, j, opts))
		}
	}
	return m, nil
}

func executePixelParallelAlgorithm(opts Opts) (img image.Image, err error) {
	m := createPNG(opts)

	wg := sync.WaitGroup{}
	wg.Add(opts.Width * opts.Height)
	for i := 0; i < opts.Width; i++ {
		for j := 0; j < opts.Height; j++ {
			go func(i, j int) {
				m.Set(i, j, getMandelbrotColor(i, j, opts))
				wg.Done()
			}(i, j)
		}
	}
	wg.Wait()
	return m, nil
}

func executeColumnParallelAlgorithm(opts Opts) (img image.Image, err error) {
	m := createPNG(opts)

	wg := sync.WaitGroup{}
	wg.Add(opts.Width)
	for i := 0; i < opts.Width; i++ {
		go func(i int) {
			for j := 0; j < opts.Height; j++ {
				m.Set(i, j, getMandelbrotColor(i, j, opts))

			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return m, nil
}

func executeWorkersAlgorithm(opts Opts) (img image.Image, err error) {
	m := createPNG(opts)

	wg := sync.WaitGroup{}
	c := make(chan map[string]int)

	for n := 0; n < 8; n++ {
		wg.Add(1)
		go func() {
			for px := range c {
				x := px["x"]
				y := px["y"]
				m.Set(x, y, getMandelbrotColor(x, y, opts))
			}
			wg.Done()
		}()
	}

	for i := 0; i < opts.Width; i++ {
		for j := 0; j < opts.Height; j++ {
			c <- map[string]int{"x": i, "y": j}
		}
	}
	close(c)
	wg.Wait()
	return m, nil
}

func executeBufferedWorkersAlgorithm(opts Opts) (img image.Image, err error) {
	m := createPNG(opts)

	wg := sync.WaitGroup{}
	c := make(chan map[string]int, opts.Width*opts.Height)
	for n := 0; n < 8; n++ {
		wg.Add(1)
		go func() {
			for px := range c {
				x := px["x"]
				y := px["y"]
				m.Set(x, y, getMandelbrotColor(x, y, opts))
			}
			wg.Done()
		}()
	}
	for i := 0; i < opts.Width; i++ {
		for j := 0; j < opts.Height; j++ {
			c <- map[string]int{"x": i, "y": j}
		}
	}
	close(c)
	wg.Wait()
	return m, nil
}

func executeBufferedColumnWorkersAlgorithm(opts Opts) (img image.Image, err error) {
	m := createPNG(opts)

	wg := sync.WaitGroup{}
	c := make(chan int, opts.Width)
	for n := 0; n < 8; n++ {
		wg.Add(1)
		go func() {
			for i := range c {
				for j := 0; j < opts.Height; j++ {
					m.Set(i, j, getMandelbrotColor(i, j, opts))
				}
			}
			wg.Done()
		}()
	}

	for i := 0; i < opts.Width; i++ {
		c <- i
	}

	close(c)
	wg.Wait()
	return m, nil
}

func createPNG(opts Opts) (m *image.RGBA) {
	palette := []color.RGBA{}
	for i := 0; i < 1000; i++ {
		c := color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255,
		}
		palette = append(palette, c)
	}
	r := image.Rect(0, 0, opts.Width, opts.Height)
	m = image.NewRGBA(r)
	return m
}

func transformColor(i int, opts Opts) color.RGBA {
	c := (float64(i) / float64(opts.MaxIteration-1)) * (255) * 15

	// if you are in set be black
	if i == opts.MaxIteration {
		return color.RGBA{
			R: uint8(0),
			G: uint8(0),
			B: uint8(0),
			A: uint8(255),
		}
	}

	// if you are in first half approach red from black
	if i < opts.MaxIteration/2-1 {
		return color.RGBA{
			R: uint8(c),
			G: uint8(0),
			B: uint8(0),
			A: uint8(255),
		}
	}

	// if you are in the second half approach white from red
	return color.RGBA{
		R: uint8(255),
		G: uint8(c),
		B: uint8(c),
		A: uint8(255),
	}
}

func getMandelbrotColor(i, j int, opts Opts) color.RGBA {
	iteration := 0
	cx := 1.5*(float64(i)-float64(opts.Width)/2.0)/(.5*float64(opts.Width)*opts.Zoom) + opts.MoveX
	cy := (float64(j)-float64(opts.Height)/2.0)/(0.5*opts.Zoom*float64(opts.Height)) + opts.MoveY

	newX := float64(0)
	newY := float64(0)
	oldX := float64(0)
	oldY := float64(0)

	for ((newX*newX)+(newY*newY) < float64(opts.Complexity)) && (iteration < opts.MaxIteration) {
		oldX = newX
		oldY = newY
		newX = oldX*oldX - oldY*oldY + cx
		newY = 2.0*oldX*oldY + cy
		iteration++
	}
	return transformColor(iteration, opts)
}

func exit(status int) {
	log.Flush()
	os.Exit(status)
}

func isUsage(err error) bool {
	return strings.HasPrefix(err.Error(), "Usage:")
}
