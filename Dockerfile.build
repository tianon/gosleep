FROM golang:1.11-alpine

RUN apk add --no-cache file git

# disable CGO for ALL THE THINGS (to help ensure no libc)
ENV CGO_ENABLED 0

ENV BUILD_FLAGS="-v -ldflags '-d -s -w'"

RUN mkdir -pv /artifacts

WORKDIR /usr/src/gosleep
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# gosleep-$(dpkg --print-architecture | awk -F- '{ print $NF }')
RUN set -x \
	&& eval "GOARCH=amd64 go build $BUILD_FLAGS -o /artifacts/gosleep-amd64" ./... \
	&& file /artifacts/gosleep-amd64 \
	&& /artifacts/gosleep-amd64 --help \
	&& time /artifacts/gosleep-amd64 --for 1s
RUN set -x \
	&& eval "GOARCH=386 go build $BUILD_FLAGS -o /artifacts/gosleep-i386" ./... \
	&& file /artifacts/gosleep-i386 \
	&& time /artifacts/gosleep-i386 --for 1s
RUN set -x \
	&& eval "GOARCH=arm GOARM=5 go build $BUILD_FLAGS -o /artifacts/gosleep-armel" ./... \
	&& file /artifacts/gosleep-armel
RUN set -x \
	&& eval "GOARCH=arm GOARM=6 go build $BUILD_FLAGS -o /artifacts/gosleep-armhf" ./... \
	&& file /artifacts/gosleep-armhf
# boo Raspberry Pi, making life hard
#RUN set -x \
#	&& eval "GOARCH=arm GOARM=7 go build $BUILD_FLAGS -o /artifacts/gosleep-armhf" ./... \
#	&& file /artifacts/gosleep-armhf
RUN set -x \
	&& eval "GOARCH=arm64 go build $BUILD_FLAGS -o /artifacts/gosleep-arm64" ./... \
	&& file /artifacts/gosleep-arm64
RUN set -x \
	&& eval "GOARCH=ppc64 go build $BUILD_FLAGS -o /artifacts/gosleep-ppc64" ./... \
	&& file /artifacts/gosleep-ppc64
RUN set -x \
	&& eval "GOARCH=ppc64le go build $BUILD_FLAGS -o /artifacts/gosleep-ppc64el" ./... \
	&& file /artifacts/gosleep-ppc64el
RUN set -x \
	&& eval "GOARCH=s390x go build $BUILD_FLAGS -o /artifacts/gosleep-s390x" ./... \
	&& file /artifacts/gosleep-s390x

RUN set -ex; \
	ls -lAFh /artifacts; \
	file /artifacts/*
