FROM golang:1.16-alpine as builder

RUN mkdir -p PickupStats/bin/
WORKDIR PickupStats/
RUN apk add --no-cache make

COPY go.mod .
COPY go.sum .
COPY Makefile .
COPY app/ ./app
COPY pkg/ ./pkg
COPY src/ ./src

RUN make app

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/PickupStats/bin/app .
COPY --from=builder /go/PickupStats/src ./src
COPY config.yaml .

EXPOSE 1323

ENTRYPOINT ["./app"]