FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o analytics-service ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/analytics-service .
COPY --from=build /app/db ./db
EXPOSE 8085
CMD ["./analytics-service"]
