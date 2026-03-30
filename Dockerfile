# -- Stage 1: Build ---------------------------------------------------
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Cache dependencies separately from source code
COPY go.mod go.sum ./
RUN go mod download

# Build static binary (no CGO = truly portable)
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o weather .


# -- Stage 2: Final image ------------------------------------------
FROM alpine:3.19

# ca-certificates: needed for HTTPS calls to OpenWeatherMap
# tzdata: needed for correct local time formating
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/weather .

ENTRYPOINT ["./weather"]
CMD ["--help"]