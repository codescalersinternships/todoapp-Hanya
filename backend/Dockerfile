FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -o server ./cmd/server.go

FROM alpine:3

WORKDIR /app

COPY --from=build /app/server .
COPY todo.db .

EXPOSE 3000

ENTRYPOINT ["./server"]

CMD ["-f", "./todo.db", "-p", "3000"]