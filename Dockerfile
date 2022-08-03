################################
# STEP 1 build executable binary
################################
FROM golang:1.17-alpine AS builder

RUN apk update && apk add --no-cache git make build-base

WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY internal internal
COPY main.go main.go

# Build the binary
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -a -installsuffix cgo -o output/app .

############################
# STEP 2 build a small image
############################
FROM alpine:3

WORKDIR /app
RUN apk update && apk add --no-cache sqlite

# Copy executable
COPY --from=builder /build/output/app /app/app
COPY configs/config.example.yml /app/config.yml
COPY views /app/views

# Expose port and declare volume
EXPOSE 8080
VOLUME /app/data

# Run the binary
ENTRYPOINT ["/app/app"]
