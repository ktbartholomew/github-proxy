FROM golang:1.9 as builder
WORKDIR /go/src/github.com/ktbartholomew/github-proxy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o github-proxy .

FROM alpine:latest
EXPOSE 8080
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/ktbartholomew/github-proxy/ .
CMD ["./github-proxy"]
