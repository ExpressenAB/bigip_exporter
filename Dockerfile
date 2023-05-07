FROM golang:1.18 AS builder

ENV GOPATH /go
ENV APPPATH /repo
COPY . /repo
RUN cd /repo && CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /repo/bigip_exporter /bigip_exporter
EXPOSE 9142
ENTRYPOINT ["/bigip_exporter"]
