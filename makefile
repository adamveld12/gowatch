dev: build
	./gowatch -debug -output="test" -args="-halt" ./test

lint: 
	find -f ./**/*.go | xargs -n 1 golint

build: clean
	go build ./cli/gowatch

clean:
	rm -rf ./gowatch
	rm -rf ./coverage.out

tests:
	go test -v -cover

test_cover:
	go test -coverprofile coverage.out
	go tool cover -html coverage.out

release:
	cd ./bin
	gox github.com/adamveld12/gowatch/cli/gowatch
