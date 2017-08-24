default: dependencies build

dependencies:
	glide install

build:
	GOPATH=$$HOME go build

build-linux:
	GOPATH=$$HOME GOOS=linux GOARCH=amd64 go build -o deployer-linux-amd64
