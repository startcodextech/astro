FROM golang:latest as dev

RUN apt update -y && apt upgrade -y
RUN apt install net-tools iproute2 -y

RUN curl -O https://raw.githubusercontent.com/angristan/openvpn-install/master/openvpn-install.sh
RUN chmod +x openvpn-install.sh

WORKDIR /go/src/github.com/start-codex/astro

RUN go install github.com/cosmtrek/air@latest

COPY . .

CMD air -c .air.toml

FROM nginx:stable as proxy

RUN apt update && apt upgrade -y

RUN apt-get install net-tools iproute2 -y