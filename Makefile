all: test build

install:
	go install

test:
	go test ./...

build:
	go build -v -ldflags="-w -s" -o output/app

docker-build:
	docker image build -t capcom6/gomvn .

docker-run: docker-build
	docker run --rm -it -p 8080:8080 capcom6/gomvn

api-docs:
	swag fmt \
	&& swag init -o ./api

view-docs:
	php -S 127.0.0.1:8080 -t ./api

clean:
	rm -r ./tmp

run:
	go run .

.PHONY: all clean install docker-build docker-run run api-docs view-docs
