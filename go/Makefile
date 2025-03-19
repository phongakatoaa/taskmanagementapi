_test:
	go test ./...

_build:
	mkdir -p "bin" && go build -o bin/

_image:
	docker build -t siransbach/task-management-api:latest .

build: _test _build

image: _test _image
