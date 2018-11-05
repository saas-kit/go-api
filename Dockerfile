# Image with necessary dependencies
FROM golang:alpine AS container
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl ca-certificates
ARG MODULE_PATH=/go/src/go-saas-kit
WORKDIR $MODULE_PATH
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download


# Go application builder
FROM container AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o /app .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o /healthchecker ./healthcheck


# Run go application
FROM scratch
ENV API_PORT=8080
WORKDIR /
COPY --from=builder /app /
COPY --from=builder /healthchecker /
COPY .env /
EXPOSE $API_PORT
CMD ["/app"]
