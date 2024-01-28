FROM golang:1.19-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/owl_clerk_bot

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
#RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go mod download && CGO_ENABLED=0 go build -o ./dist/owl_clerk_bot cmd/app/main.go

# Start fresh from a smaller image
FROM alpine:3.18.6
RUN apk add ca-certificates

COPY --from=build_base /tmp/owl_clerk_bot/dist/owl_clerk_bot /app/owl_clerk_bot

# This container exposes port 8080 to the outside world
EXPOSE 8000

# Run the binary program produced by `go install`
CMD ["/app/owl_clerk_bot"]
