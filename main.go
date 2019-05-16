package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	//"errors"
	"github.com/joho/godotenv"
	"gopkg.in/h2non/bimg.v1"
	"path/filepath"
)

const version = "0.0.1"

var options struct {
	version  bool
	dstWidth int
	input    string
	output   string
	format   string
	convert  string
	medium   string
}

func main() {
	var originalBuf []byte
	var newBuf []byte
	var err error

	// we don't handle errors on purpose, the system should accept both
	// .env file and traditional environment variables
	godotenv.Load()

	flag.BoolVar(&options.version, "-version", false, "Check version")
	flag.StringVar(&options.input, "-image", "", "Input image file")
	flag.IntVar(&options.dstWidth, "-width", 0, "Output width dimension")
	flag.StringVar(&options.output, "-output", "", "Output directory")
	flag.StringVar(&options.convert, "-convert", "", "Convert to supported type: jpeg,png, webp,tiff,gif,pdf,svg")
	flag.StringVar(&options.medium, "-medium", "", "Choose medium: disk or s3, default is disk")
	flag.BoolVar(&options.version, "v", false, "Check version")
	flag.StringVar(&options.input, "i", "", "Input image file")
	flag.IntVar(&options.dstWidth, "w", 0, "Output width dimension")
	flag.StringVar(&options.output, "o", "", "Output directory")
	flag.StringVar(&options.convert, "c", "", "Convert to supported type: jpeg,png, webp,tiff,gif,pdf,svg")
	flag.StringVar(&options.medium, "m", "", "Choose medium: disk or s3, default is disk")
	flag.Parse()

	if options.version {
		fmt.Println("Version " + version)
		os.Exit(0)
	}

	if options.input == "" {
		log.Fatal("No input file")
		os.Exit(1)
	}

	originalBuf, err = fileToBuf(options.input)
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

	if options.medium == "" {
		options.medium = "disk"
	}

	newBuf, err = Resize(originalBuf, options.dstWidth)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not resize the image")
		os.Exit(1)
	}

	newBuf, err = Convert(newBuf, bimg.PNG)

	if err != nil {
		log.Fatal(err)
		log.Fatal("Could Convert image")
	}

	err = SaveImg(options.medium, options.output, options.input, newBuf)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not save the file to " + options.medium)
	}

	if len(flag.Args()) == 0 {
		os.Exit(1)
	}
}

// Reads a file into buffer
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

// Uploads a file to s3
func saveTos3(fullpath string, buf []byte) error {
	return nil
}

// Saves a bimg *Image to defined medium and path.
func SaveImg(medium string, dstPath string, inPath string, buf []byte) error {
	var err error
	err = nil
	filename, _ := getFilename(inPath)
	typename := bimg.DetermineImageTypeName(buf)
	fullpath := dstPath + "/" + filename + "." + typename
	if medium == "s3" || medium == "S3" {
		err = saveTos3(fullpath, buf)
	}
	if medium == "disk" {
		err = bimg.Write(fullpath, buf)
	}
	return err
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
