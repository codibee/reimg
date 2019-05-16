package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"gopkg.in/h2non/bimg.v1"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const version = "0.2.1"

var options struct {
	version  bool
	dstWidth int
	covtype  int
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
	if options.convert != "" {
		convtype, err := convertTo(options.convert)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		newBuf, err = Convert(newBuf, convtype)
		if err != nil {
			log.Fatal(err)
			log.Fatal("Could Convert image")
			os.Exit(1)
		}
	}

	err = SaveImg(options.medium, options.output, options.input, newBuf)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not save the file to " + options.medium)
	}
	fmt.Println("File successfully saved to " + options.medium)
	if len(flag.Args()) == 0 {
		os.Exit(1)
	}
}

// Determines if we need image convertion or not
func convertTo(imageType string) (bimg.ImageType, error) {

	var convtype bimg.ImageType
	var err error
	convtype = -1
	err = nil

	if imageType == "jpg" {
		imageType = "jpeg"
	}

	for k, v := range bimg.ImageTypes {
		if v == imageType {
			convtype = k
		}
	}

	if convtype < 0 {
		err = errors.New("Cannot convert to an unsupported type " + imageType)
	}
	return convtype, err
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

	var err error
	err = nil

	if os.Getenv("AWS_KEY") == "" || os.Getenv("AWS_SECRET") == "" ||
		os.Getenv("AWS_BUCKET") == "" || os.Getenv("AWS_REGION") == "" {
		return errors.New("Missing AWS credentials, please use .env file or real environment variables")
	}

	s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	if err != nil {
		return err
	}

	err = s3Add(s, fullpath, buf)

	return err
}

// Adds object to S3 via defined aws session
func s3Add(s *session.Session, fullpath string, buf []byte) error {
	var size int64
	var acl string
	size = int64(len(buf))
	acl = "private"

	if os.Getenv("AWS_ACL") != "" {
		acl = os.Getenv("AWS_ACL")
	}
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(os.Getenv("AWS_BUCKET")),
		Key:                aws.String(fullpath),
		ACL:                aws.String(acl),
		Body:               bytes.NewReader(buf),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buf)),
		ContentDisposition: aws.String("attachment"),
	})

	return err
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
	} else if medium == "disk" {
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			os.Mkdir(dstPath, 0755)
		}
		if err == nil {
			err = bimg.Write(fullpath, buf)
		}
	} else {
		err = errors.New("Unsupported medium " + medium)
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
