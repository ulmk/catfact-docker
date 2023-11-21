FROM golang:latest

COPY ./ ./

ENV GOPATH=""

RUN go build -o main .

CMD ["./main"]