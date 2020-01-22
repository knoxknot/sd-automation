FROM golang:1.12-stretch AS base

WORKDIR /src/
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 

COPY . /src/

RUN go mod download
RUN go build -o /src/boardingapi

FROM scratch

WORKDIR /src/
COPY --from=base /src/boardingapi /src/boardingapi

ENTRYPOINT ["/src/boardingapi"]