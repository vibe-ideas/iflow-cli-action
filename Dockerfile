FROM golang:1.22-alpine

WORKDIR /app

COPY main.go .

RUN go build -o iflow-action main.go

ENTRYPOINT ["/app/iflow-action"]
