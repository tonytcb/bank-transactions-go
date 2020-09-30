FROM golang:1.14-stretch

MAINTAINER Tony C. Batista

ENV TIMEZONE America/Sao_Paulo

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod vendor
RUN go get github.com/pilu/fresh

EXPOSE 8080
