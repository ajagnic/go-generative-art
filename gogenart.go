package main

import (
	"flag"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/ajagnic/go-generative-art/sketch"
)

func main() {
	i := flag.Int("i", 10000, "number of iterations")
	w := flag.Uint("width", 2400, "desired width of image")
	h := flag.Uint("height", 1600, "desired height of image")
	shake := flag.Float64("shake", 0.0, "amount to randomize pixel positions")
	min := flag.Uint("min", 3, "minimum number of polygon sides")
	max := flag.Uint("max", 5, "maximum number of polygon sides")
	fill := flag.Int("fill", 1, "1 in N chance to fill polygon")
	color := flag.Int("color", 0, "1 in N chance to randomize polygon color")
	s := flag.Float64("s", 0.1, "polygon size (percentage of width)")
	output := flag.String("o", "", "file to use as output")
	flag.Parse()

	in := handleInput()
	defer in.Close()

	img, enc, err := image.Decode(in)
	if err != nil {
		log.Fatalf("could not decode: %v\n", err)
	}

	if *max < *min {
		min, max = max, min
	}
	canvas := sketch.NewSketch(img, sketch.Params{
		Iterations:         *i,
		Width:              int(*w),
		Height:             int(*h),
		PixelShake:         int(*shake * float64(*w)),
		PolygonSidesMin:    int(*min),
		PolygonSidesMax:    int(*max),
		PolygonFillChance:  *fill,
		PolygonColorChance: *color,
		PolygonSizeRatio:   *s,
	})
	canvas.Draw()

	out, enc := handleOutput(*output, enc)
	defer out.Close()

	switch enc {
	case "png":
		png.Encode(out, canvas.Image())
	default:
		jpeg.Encode(out, canvas.Image(), nil)
	}
}

func handleInput() (in *os.File) {
	var err error
	if args := flag.Args(); len(args) > 0 {
		in, err = os.Open(args[0])
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		in = os.Stdin
	}
	return
}

func handleOutput(file, enc string) (*os.File, string) {
	if file != "" {
		out, err := os.Create(file)
		if err != nil {
			log.Fatalln(err)
		}
		fSlc := strings.Split(file, ".")
		return out, fSlc[len(fSlc)-1]
	} else {
		return os.Stdout, enc
	}
}
