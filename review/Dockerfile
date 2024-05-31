# build stage
FROM golang:alpine3.20 as builder

WORKDIR /go/src/coupon-service

COPY . .

WORKDIR /go/src/coupon-service

RUN mkdir /opt/service

RUN go build -o /opt/service/main ./cmd/coupon_service/main.go

RUN chmod +x /opt/service/main

WORKDIR /opt/service

COPY config.yml /opt/service/config.yml

EXPOSE 8080

ENTRYPOINT ["./main"]
