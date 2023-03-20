FROM golang:1.18

EXPOSE 80

WORKDIR /go/src/views-service
COPY ./ ./
RUN go mod download
RUN go build -o views-service

CMD ["./views-service"]
