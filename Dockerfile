ARG GOLANG_VERSION=1.19-alpine

FROM golang:${GOLANG_VERSION} AS build
WORKDIR /build
COPY . .

RUN go mod vendor

RUN go build -o /bin/server -mod=vendor

FROM alpine:latest AS dev

COPY --from=build /bin/server /bin/server

EXPOSE 44444:44444
WORKDIR /bin
ENTRYPOINT ["/bin/server"]
