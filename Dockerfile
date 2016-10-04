FROM alpine:latest

EXPOSE 9142

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/ExpressenAB/bigip_exporter
COPY . $APPPATH
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc \
    && cd $APPPATH && go get -u github.com/kardianos/govendor && $GOPATH/bin/govendor build +p \
    && cp bigip_exporter /bigip_exporter && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT ["/bigip_exporter"]
