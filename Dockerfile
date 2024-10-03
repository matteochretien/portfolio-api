# Dockerfile to containerize a go application, the first stage is a build stage and second is the runtime stage that will contains a healthcheck

# Build stage
FROM golang:1.22.5-alpine3.20 as builder

# Maintainer
LABEL maintainer="Mattéo Chrétien <contact@matteochretien.com>"

RUN apk add git gcc build-base

# Set the Current Working Directory inside the container
WORKDIR /app

ENV CGO_ENABLED=1
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

# Copy go mod and sum files
COPY go.mod go.sum ./
# download dependencies and cache them using buildkit
RUN --mount=type=cache,target=/gomod-cache go mod download

COPY . .
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache go build -ldflags='-s -w' -o /dist/app cmd/main.go
RUN chmod +x /dist/app

# Runtime stage, the binary will run in a user mode
FROM alpine:3.20
# Maintainer
LABEL maintainer="Mattéo Chrétien <contact@matteochretien.com>"

# Install curl
RUN apk update && apk add --no-cache curl

# Set the Current Working Directory inside the container
WORKDIR /app

# Create user appuser to run the application
RUN adduser -D -g '' appuser
# Switch to non-root user
USER appuser

# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chown=appuser:appuser /dist/app /app/app

# Expose port 3000 to the outside world
EXPOSE 3000
# Command to run the executable
CMD ["/app/app"]
# Healthcheck
HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 CMD curl -f http://localhost:3000/health || exit 1