FROM golang:1.10.2-alpine as builder
ENV CGO_ENABLED=0
RUN apk --no-cache add git make upx curl
WORKDIR $GOPATH/src/github.com/stamm/dep_radar
COPY . $GOPATH/src/github.com/stamm/dep_radar
RUN make build && \
	upx $GOPATH/bin/dep_radar  && \
	cp $GOPATH/bin/dep_radar /


FROM alpine:3.7
MAINTAINER Rustam Zagirov <stammru@gmail.com>
RUN apk --no-cache add ca-certificates
EXPOSE 8081
VOLUME ["/cfg/"]
ENTRYPOINT ["/bin/dep_radar"]
COPY --from=builder /dep_radar /bin/

