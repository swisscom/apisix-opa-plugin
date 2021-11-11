FROM golang:1.17-alpine
RUN apk add --no-cache make
WORKDIR /app
COPY / /app
RUN make