FROM golang:1.20.0-alpine AS build

RUN apk update && apk upgrade && apk add --no-cache bash git

RUN mkdir -p /app/receipt-processor-challenge
WORKDIR /app/receipt-processor-challenge

COPY go.mod  .
RUN go mod download
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o app cmd/receipt-processor-challenge/main.go

FROM alpine:latest as prod
RUN apk --no-cache add ca-certificates
COPY --from=build /app/receipt-processor-challenge/app /app
EXPOSE 8080
CMD ["/app"]
