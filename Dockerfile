FROM golang:1.17.3-alpine3.14 as builder
RUN apk add git
RUN mkdir /reservations-service
ADD . /reservations-service
WORKDIR /reservations-service

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

RUN mkdir /reservations-service

WORKDIR /reservations-service/

COPY --from=builder /reservations-service/main .

ARG DBpw_arg=default_value 
ENV DBpw=$DBpw_arg

EXPOSE 8080

CMD ["./main"]