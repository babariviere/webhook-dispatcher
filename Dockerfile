FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -v -o /webhook-dispatcher main.go

FROM alpine
RUN apk add ca-certificates

COPY --from=build /webhook-dispatcher /webhook-dispatcher

EXPOSE 8080

CMD ["/webhook-dispatcher"]
