clean:
	@rm -rf dist
	@mkdir -p dist

build: clean
	GOOS=linux go build -o dist/reimg

requirements:
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/aws/session
	go get github.com/aws/aws-sdk-go/service/s3
	go get github.com/joho/godotenv
	go get gopkg.in/h2non/bimg.v1

install:
	@go install
	@echo "You can find your bin at ${GOPATH}/bin"
