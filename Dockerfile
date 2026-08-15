# ---------- Stage 1: Build ----------
FROM golang:1.26-alpine AS build

WORKDIR /src

# Download dependencies first
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build a static Linux binary
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" -o /app ./cmd


# ---------- Stage 2: Runtime ----------
FROM scratch

# Copy only the compiled binary
COPY --from=build /app /app

ENTRYPOINT ["/app"]