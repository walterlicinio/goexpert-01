FROM golang:1.20

ENV CGO_ENABLED=1
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o server

CMD [ "./server" ]
