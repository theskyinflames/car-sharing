FROM golang:1.19.3-alpine3.16

WORKDIR /challenge

COPY . .

RUN go build cmd/main.go cmd/run.go

FROM alpine:3.14.0
WORKDIR /callenge

COPY --from=0 /challenge/main .
EXPOSE 80
ENTRYPOINT [ "./main" ]



