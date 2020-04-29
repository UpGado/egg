package main

import (
	"flag"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strings"
)

func imread(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func vstack(top_img image.Image, bot_img image.Image) image.Image {
	top_width, top_height := top_img.Bounds().Max.X, top_img.Bounds().Max.Y
	bot_width, bot_height := bot_img.Bounds().Max.X, bot_img.Bounds().Max.Y
	if top_width != bot_width {
		log.Fatal("imgs need to be of same width")
	}
	img := image.NewRGBA(image.Rect(0, 0, bot_width, top_height+bot_height))
	for y := 0; y < top_height; y++ {
		for x := 0; x < top_width; x++ {
			img.Set(x, y, top_img.At(x, y))
		}
	}
	for y := 0; y < bot_height; y++ {
		for x := 0; x < bot_width; x++ {
			img.Set(x, y+top_height, bot_img.At(x, y))
		}
	}
	return img
}

func hatch(mom_egg image.Image, mom_egg_template string,
	new_egg_template string) image.Image {
	mom_width, mom_height := mom_egg.Bounds().Max.X, mom_egg.Bounds().Max.Y
	if new_egg_template == "" {
		return image.NewRGBA(image.Rect(0, 0, mom_width, 0))
	}
	quadrant_height := mom_height / len(mom_egg_template)
	quadrant_i := strings.Index(mom_egg_template, string(new_egg_template[0]))
	quadrant := image.NewRGBA(image.Rect(0, 0, mom_width, quadrant_height))
	for y := 0; y < quadrant_height; y++ {
		for x := 0; x < mom_width; x++ {
			quadrant.Set(x, y, mom_egg.At(x, (quadrant_i*quadrant_height)+y))
		}
	}
	return vstack(quadrant, hatch(mom_egg, mom_egg_template, new_egg_template[1:]))
}

func main() {
	eggpath := flag.String("src", "egg.jpg", "path for mom egg")
	new_template := flag.String("new", "eegg", "what new egg should look like")
	mom_template := flag.String("template", "egg", "template for mom egg")
	dstpath := flag.String("dst", "new_egg.jpg", "where to save new egg")
	flag.Parse()
	mom_egg := imread(*eggpath)
	new_egg := hatch(mom_egg, *mom_template, *new_template)
	file, err := os.Create(*dstpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	jpeg.Encode(file, new_egg, nil)
}
