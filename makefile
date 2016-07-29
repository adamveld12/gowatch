dev: build
	./gowatch -debug -output="test" -args="-halt" ./test

lint: 
	find -f ./**/*.go | xargs -n 1 golint

build: clean gowatch

clean:
	rm -rf ./gowatch ./coverage.out ./bin

tests:
	go test -v -cover

cover: clean coverage.out
	go tool cover -html coverage.out

release: gox ./bin/
	gox -rebuild -arch="amd64" -output="./bin/{{.Dir}}_{{.OS}}" github.com/adamveld12/gowatch/cli/gowatch

.PHONY: cover test clean lint dev release

./bin/:
	mkdir ./bin

gox:
	go get github.com/mitchellh/gox

gowatch:
	go build ./cli/gowatch

coverage.out:
	go test -coverprofile coverage.out

