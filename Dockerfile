FROM golang AS build-env

RUN GO111MODULE=off go get -u github.com/esrrhs/fastreplace
RUN GO111MODULE=off go get -u github.com/esrrhs/fastreplace/...
RUN GO111MODULE=off go install github.com/esrrhs/fastreplace

FROM debian
COPY --from=build-env /go/bin/fastreplace .
WORKDIR ./
