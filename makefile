dev: build
	./gowatch -debug -output="test" -args="-halt" ./test

lint: 
	find -f ./**/*.go | xargs -n 1 golint

build: clean gowatch

clean:
	rm -rf ./gowatch ./coverage.out

tests:
	go test -v -cover

cover: clean coverage.out
	go tool cover -html coverage.out

coverage.out:
	go test -coverprofile coverage.out

release:
	cd ./bin
	gox github.com/adamveld12/gowatch/cli/gowatch

.PHONY: cover test clean lint dev

gowatch:
	go build ./cli/gowatch
