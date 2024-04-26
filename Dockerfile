#Download base image alpine 3.10
FROM golang:1.22.1-alpine3.19 AS build-env

RUN apk update && apk add --no-cache binutils git curl unzip tar

ENV WORKDIR /go/src/github.com/zalf-rpm/Hermes2Go
WORKDIR ${WORKDIR}
COPY . ${WORKDIR}

RUN git describe --always --tags --long > /version.txt
WORKDIR /go/src/github.com/zalf-rpm/Hermes2Go/src/hermes2go
RUN go get gopkg.in/yaml.v3
RUN VERSION=$(cat /version.txt) && go build -v -ldflags "-X main.version=${VERSION}"
WORKDIR /go/src/github.com/zalf-rpm/Hermes2Go/src/calcHermesBatch
RUN VERSION=$(cat /version.txt) && go build -v -ldflags "-X main.version=${VERSION}"

FROM alpine:3.19

COPY --from=build-env /go/src/github.com/zalf-rpm/Hermes2Go/src/hermes2go/hermes2go /hermes2go/
COPY --from=build-env /go/src/github.com/zalf-rpm/Hermes2Go/src/calcHermesBatch/calcHermesBatch /hermes2go/
RUN chmod -R 555 /hermes2go

ENV PATH=/hermes2go:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN addgroup -S mygroup && adduser -S myuser -G mygroup
USER myuser

CMD ["/bin/bash"]