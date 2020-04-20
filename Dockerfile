#Wirebuard & wg-control-linux WebUI
ARG COMMIT="0.0.1"

FROM golang:alpine AS build-back
RUN apk update && apk add --no-cache git
WORKDIR /app
ARG COMMIT
COPY . .
RUN go get github.com/TheLazarusNetwork/wg-control
RUN go build -o wg-control-linux github.com/TheLazarusNetwork/wg-control

FROM node:10-alpine AS build-front
WORKDIR /app
RUN apk update && apk add --no-cache git
RUN git clone https://github.com/TheLazarusNetwork/wg-control.git
WORKDIR /app/wg-control/ui
RUN npm install
RUN npm run build

FROM alpine
WORKDIR /app
COPY --from=build-back /app/wg-control-linux .
COPY --from=build-front /app/wg-control/ui/dist ./ui/dist
COPY .env .
RUN chmod +x ./wg-control-linux
RUN apk add --no-cache ca-certificates
RUN apk update && apk add --no-cache wireguard-tools
RUN mkdir /etc/wireguard/keys
EXPOSE 9080

CMD ["/app/wg-control-linux"]
