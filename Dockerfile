FROM golang:1.14-stretch

WORKDIR /app

COPY . .

RUN go mod download
RUN go get github.com/pilu/fresh

EXPOSE 8080