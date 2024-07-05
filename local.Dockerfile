# Step 1: Modules caching
FROM golang:1.22.5-alpine3.19 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22.5-alpine3.19 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go build -o /bin/app ./cmd/app

# Step 3: Deploy
FROM alpine:3.19.1
WORKDIR /orgonaut
COPY --from=builder /app/configs /orgonaut/configs
COPY --from=builder /bin/app /orgonaut/
CMD ["/orgonaut/app"]