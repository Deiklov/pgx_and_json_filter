FROM golang:1.16-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./cronjob ./cronjob

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates
COPY --from=build_base /app/cronjob /app/cronjob
WORKDIR /app


# Run the binary program produced by `go install`
CMD ["./cronjob"]