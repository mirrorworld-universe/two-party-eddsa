ARG GO_VERSION

#FROM golang:${GO_VERSION} as builder
FROM golang:1.18-alpine as builder

WORKDIR /builder/

COPY . .

ENV PATH="/go/bin:${PATH}"
ENV TZ=Asia/Shanghai \
    CGO_ENABLED=1   \
    GO111MODULE=on  \
    GOOS=linux
#    GOPROXY=https://goproxy.cn,direct

#RUN go get ./...
RUN go mod download
RUN go build -tags musl --ldflags "-extldflags -static" -o mirror-market-backend .

FROM scratch

ENV TZ=Asia/Shanghai \
    GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    PROGRAM_ENV=pro

WORKDIR /app

COPY --from=builder /builder/mirror-market-backend .
COPY ./conf ./conf

# for http
EXPOSE 3000

# 启动服务
CMD ["./two-party-eddsa"]