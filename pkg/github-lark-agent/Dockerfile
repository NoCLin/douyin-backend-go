FROM golang:1.18-alpine as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./*.go ./

RUN go build -o ./agent


FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/agent .

EXPOSE 8081

CMD ["./agent"]


