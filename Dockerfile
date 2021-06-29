FROM alpine:latest

EXPOSE 9142

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/ExpressenAB/bigip_exporter
COPY . $APPPATH
WORKDIR $APPPATH
RUN apk --no-cache add wget ca-certificates go libstdc++ govendor
RUN apk add --update -t build-deps go mercurial libc-dev gcc libgcc
RUN go mod init && go mod vendor
RUN govendor build +p && cp bigip_exporter /bigip_exporter && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT ["/bigip_exporter"]
