all: test build

install:
	go install

test:
	go test ./...

build:
	go build -v -ldflags="-w -s" -o output/app

.PHONY: docker-build
docker-build:
	docker image build -t capcom6/gomvn .

.PHONY: docker-run
docker-run: docker-build
	docker run --rm -it -p 8080:8080 capcom6/gomvn

clean:
	rm -r output

.PHONY: run
run:
	go run .
