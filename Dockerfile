FROM golang:1.21.0-alpine3.18 as build
RUN apk add --no-cache gcc libc-dev ca-certificates && update-ca-certificates
WORKDIR /app

ENV CGO_ENABLED=0
ENV GO111MODULE=on

COPY go.mod ./
RUN go mod download
COPY . .

RUN go build -o /app/main .

FROM alpine:3.18 AS final
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main /app
RUN mkdir /app/static
EXPOSE 8088
CMD ["/app/main", "-l", "8088", "-r", "/app/static"]

