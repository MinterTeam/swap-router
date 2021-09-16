FROM golang:1.15.0-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers
WORKDIR /app
COPY . /app
RUN go build -o ./builds/app ./app/app.go

FROM alpine:3.7

COPY --from=builder /app/builds/app /usr/bin/swaprouter
RUN addgroup minteruser && adduser -D -h /minter -G minteruser minteruser
USER minteruser
WORKDIR /minter
CMD ["/usr/bin/swaprouter"]
