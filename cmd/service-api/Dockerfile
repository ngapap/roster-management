FROM golang:alpine as stage-builder

RUN apk update && apk upgrade && apk add --no-cache bash git

ENV GOFLAGS=-mod=vendor

WORKDIR /app

# Copy project files
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY migrations ./migrations
COPY cmd/service-api ./cmd/service-api
COPY configs ./configs
COPY internal ./internal
COPY pkg ./pkg

# Build binaries
RUN GOOS=linux GOARCH=amd64 go build -o /app/roster /app/cmd/service-api
RUN GOOS=linux GOARCH=amd64 go build -o /app/migrate /app/migrations

# Make sure binaries are executable
RUN chmod +x /app/migrate /app/roster

# Use shell form for chaining commands
ENTRYPOINT sh -c "/app/migrate up && /app/roster"
