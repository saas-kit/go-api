# Image with necessary dependencies
FROM golang:alpine AS container
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl ca-certificates
WORKDIR /go/src/go-saas-kit
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
WORKDIR /
COPY --from=builder /app /
COPY --from=builder /healthchecker /
COPY .env /
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=1s --start-period=100s --retries=3 CMD [ "/healthchecker" ]
CMD ["/app"]
