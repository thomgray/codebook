run:
	codebookdevmode=true go run .

build:
	go build -o target/codebook .code

install:
	go install

test:
	go test ./...