FROM golang:1.21

WORKDIR /app

COPY . .

RUN go build -o service_b .

EXPOSE 8081

CMD [ "./service_b"]