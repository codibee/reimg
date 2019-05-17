# reimg

Reimg is a simple application uses libvips in order to resize images and saves them to disk or aws s3. It is based on [bimg](https://github.com/h2non/bimg) and the main goal is to use it in order to resize, convert etc and then upload the results into S3. 

The application will grow over time supporting more features for images including a lambda function version. 

## Installation

#### Makefile
```bash
# do fetch the required dependencies
make requirements

# install the application to your $GOPATH/bin folder (see bellow)
make install

# or build the application in place
make build

# if you need to clean your builds
make clean
```

## Usage
```bash
  --convert string
    	Convert to supported types: jpeg,png, webp,tiff,gif,pdf,svg
  --image string
    	Input image file
  --medium string
    	Choose medium: disk or s3, default is disk
  --output string
    	Output directory
  --version
    	Check version
  --width int
    	Output width dimension

# The example bellow resizes the image sample.jpg to 1024 px width, converts it to jpg 
# and finally saves it to disk (default option) inside the folder "resized"

reimg --image sample.png --output resized --width 1024 --convert jpg

# The example bellow resizes the image sample.jpg to 1024 px width, converts it to jpg 
# and finally saves it to s3, inside the folder "resized". In order to use s3 save you 
# need to have a .env file or real environment variables for AWS credentials. 
# please check .env.txt

reimg --image sample.png --output resized --medium s3 --width 1024 --convert jpg
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[Lesser GPL v3](LICENSE.md)

![LGPLv3](https://www.gnu.org/graphics/agplv3-88x31.png)
