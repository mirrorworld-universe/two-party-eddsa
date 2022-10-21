ARG GO_VERSION

FROM golang:1.18-alpine as builder

WORKDIR /builder/
ADD . /builder/

ENV TZ=Asia/Shanghai \
    GO111MODULE=on

RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o two-party-eddsa

FROM scratch

ENV TZ=Asia/Shanghai \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    PROGRAM_ENV=pro

WORKDIR /app

COPY --from=builder /builder/two-party-eddsa .
COPY ./conf ./conf

EXPOSE 3000
CMD ["./two-party-eddsa"]