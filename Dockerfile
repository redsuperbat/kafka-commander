FROM golang:alpine3.14 as compiler

WORKDIR /app/build

COPY . .

RUN go get

RUN go build

FROM alpine:3.14

WORKDIR /app/prod

COPY --from=compiler /app/build/kafka-commander .

ENV GIN_MODE=release

ENTRYPOINT ["./kafka-commander"]
