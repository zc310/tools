package main

import (
	"errors"
	"flag"
	"fmt"
	"gotest/utils/file"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/nfnt/resize"

	"github.com/artyom/smartcrop"
	"github.com/disintegration/imaging"
)

func cut(fn string) {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	var img image.Image
	if strings.ToLower(path.Ext(fn)) == ".jpg" {
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Fatal(fn, "\t", err)
		}
	} else {
		img, err = png.Decode(file)
		if err != nil {
			log.Fatal(fn, "\t", err)
		}
	}

	file.Close()

	img1 := resize.Resize(900, 0, img, resize.NearestNeighbor)

	img2 := imaging.Fill(img, 900, 500, imaging.Center, imaging.Lanczos)
	savefile(img2, strings.Split(fn, ".")[0]+"_resized")
	topCrop, err := smartcrop.Crop(img1, 900, 500)
	if err != nil {
		log.Fatal(err)
	}
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	sub, ok := img1.(SubImager)
	if ok {
		cropImage := sub.SubImage(topCrop)
		savefile(cropImage, strings.Split(fn, ".")[0]+"_smart")
	} else {
		log.Fatal(errors.New("No SubImage support"))
	}

}
func savefile(img image.Image, fn string) {
	out, err := os.Create("test_" + fn + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	png.Encode(out, img)
}

func main() {

	var fn = flag.String("f", "test.jpg", "jpg file")
	flag.Parse()
	if !file.Exists(*fn) {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}

		var ext string
		for _, file := range files {
			ext = strings.ToLower(path.Ext(file.Name()))
			if ((ext == ".jpg") || (ext == ".png")) && !(strings.HasPrefix(file.Name(), "test_")) {
				cut(file.Name())

			}
		}

		log.Fatal(errors.New(fmt.Sprintf("%s not find!", *fn)))

	}

}
