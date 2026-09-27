FROM golang:1.27.1-bookworm AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/fastreplace .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /out/fastreplace ./fastreplace
ENTRYPOINT ["./fastreplace"]
