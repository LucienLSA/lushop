FROM golang:1.21.0-alpine3.18 AS builder

WORKDIR /build

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct
ENV GOARCH amd64

COPY . .
RUN go mod tidy

RUN go build -ldflags="-s -w" -o app main.go

FROM alpine

WORKDIR /build

ENV TZ Asia/Shanghai

COPY --from=builder /build/app .
WORKDIR /
COPY config-pro.yaml /config-pro.yaml

ENV JZIN_PRO true

#EXPOSE 20015

CMD ["./build/app"]