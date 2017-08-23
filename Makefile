default: dependencies build

dependencies:
	glide install

build:
	GOPATH=$$HOME go build
