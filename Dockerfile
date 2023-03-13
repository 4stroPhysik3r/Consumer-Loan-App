FROM golang:1.17-alpine3.14

WORKDIR /consumer-loan-app
COPY . /consumer-loan-app
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
