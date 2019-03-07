
run: build
	./bin/artisan

build:
	go build -o bin/artisan cmd/*.go