package main

import (
	//	"errors"
	"flag"
	"fmt"
	"gopkg.in/h2non/bimg.v1"
	"log"
	"os"
	"path/filepath"
)

const version = "0.0.1"

var options struct {
	version  bool
	input    string
	output   string
	format   string
	dstWidth int
}

func main() {
	flag.BoolVar(&options.version, "-version", false, "Check version")
	flag.StringVar(&options.input, "-image", "", "Input image file")
	flag.IntVar(&options.dstWidth, "-width", 0, "Output width dimension")
	flag.StringVar(&options.output, "-output", "", "Output directory")
	flag.BoolVar(&options.version, "v", false, "Check version")
	flag.StringVar(&options.input, "i", "", "Input image file")
	flag.IntVar(&options.dstWidth, "w", 0, "Output width dimension")
	flag.StringVar(&options.output, "o", "", "Output path")
	flag.Parse()

	if options.version {
		fmt.Println("Version " + version)
		os.Exit(0)
	}

	if options.input == "" {
		log.Fatal("No input file")
		os.Exit(1)
	}

	imageBuf, err := fileToBuf(options.input)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not read image file")
		os.Exit(1)
	}

	if options.output == "" {
		log.Fatal("No output directory")
		os.Exit(1)
	}

	if options.dstWidth <= 0 {
		log.Fatal("Width dimention should be greater than 0")
		os.Exit(1)
	}

	imageResized, err := Resize(imageBuf, options.dstWidth)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not resize the image")
		os.Exit(1)
	}
	convBuf, err := Convert(imageResized, bimg.PNG)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Could Convert image")
	}

	err = SaveImg("disk", options.output, options.input, convBuf)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not save the file to medium")
	}

	if len(flag.Args()) == 0 {
		os.Exit(1)
	}
}

func fileToBuf(file string) ([]byte, error) {
	buffer, err := bimg.Read(file)
	if err != nil {
		return nil, err
	}
	return buffer, err
}

// Returns the filename with/without extension
func getFilename(path string) (string, string) {
	var file = filepath.Base(path)
	var ext = filepath.Ext(file)
	var name = file[0 : len(file)-len(ext)]
	return name, ext
}

// Converts a bimg *Image to defined type
func Convert(buf []byte, imgtype bimg.ImageType) ([]byte, error) {
	image := bimg.NewImage(buf)
	newBuff, err := image.Convert(imgtype)
	if err != nil {
		return nil, err
	}
	return newBuff, nil
}

// Saves a bimg *Image to defined medium and path.
func SaveImg(medium string, dstPath string, inPath string, buf []byte) error {
	filename, _ := getFilename(inPath)
	typename := bimg.DetermineImageTypeName(buf)
	fmt.Println(typename)

	if medium == "disk" {
		bimg.Write(dstPath+"/"+filename+"."+typename, buf)
		return nil
	}
	return nil
}

// Resize images based on bimg package
func Resize(buf []byte, dstWidth int) ([]byte, error) {
	sizes, err := bimg.Size(buf)
	if err != nil {
		return nil, err
	}

	aspect_ratio := float64(sizes.Width) / float64(sizes.Height)
	dstHeightF := float64(dstWidth) / aspect_ratio
	dstHeight := int(dstHeightF)

	newImage, err := bimg.NewImage(buf).Resize(dstWidth, dstHeight)
	if err != nil {
		return nil, err
	}
	return newImage, err
}
