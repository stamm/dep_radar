FROM golang:1.9.4-alpine as builder
ENV CGO_ENABLED=0
RUN apk --no-cache add wget git make && \
  mkdir -p $GOPATH/src/github.com/stamm/dep_radar
COPY . $GOPATH/src/github.com/stamm/dep_radar
WORKDIR $GOPATH/src/github.com/stamm/dep_radar
RUN make dep_install
RUN make deps
RUN make build
RUN cp $GOPATH/bin/dep_radar /


FROM alpine:3.6
MAINTAINER Rustam Zagirov <stammru@gmail.com>

RUN apk --no-cache add ca-certificates
COPY --from=builder /dep_radar /bin/
EXPOSE 8081
VOLUME ["/cfg/"]
ENTRYPOINT ["/bin/dep_radar"]

