install:
	dep ensure

build-api:
	go build -o ${GOBIN}/tmga.exe main.go

build-genetic:
	go build -o ${GOBIN}/genetic.exe genetic.go