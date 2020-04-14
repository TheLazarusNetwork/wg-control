ARG COMMIT="0.0.1"

FROM golang:alpine AS build-back
WORKDIR /app
ARG COMMIT
COPY . .
RUN go build -ldflags="-X 'github.com/TheLazarusNetwork/wg-control/util.Version=${COMMIT}'" -o wg-control-linux

FROM node:10-alpine AS build-front
WORKDIR /app
COPY ui/package*.json ./
RUN npm install
COPY ui/ ./
RUN npm run build

FROM alpine
WORKDIR /app
COPY --from=build-back /app/wg-control-linux .
COPY --from=build-front /app/dist ./ui/dist
COPY .env .
RUN chmod +x ./wg-control-linux
RUN apk add --no-cache ca-certificates
EXPOSE 9080

CMD ["/app/wg-control-linux"]