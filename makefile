install:
	dep ensure

build:
	go build -o ${GOBIN}/tmga.exe main.go
