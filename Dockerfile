FROM golang:1.15.8
WORKDIR /go/src/github.com/AttilioGreco/smtp-debugger
COPY go.sum .
COPY go.mod .
COPY server.go .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/AttilioGreco/smtp-debugger/smtp-to-console .
CMD ["./smtp-to-console"]