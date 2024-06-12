FROM golang:1.22-alpine as build

# update index packages for OS
RUN apk update

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod tidy

COPY . .

RUN go version

RUN go build -o product .

FROM alpine:latest

WORKDIR /root/

RUN apk update && apk add tzdata

RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

RUN echo "Asia/Jakarta" > /etc/timezone
RUN apk del tzdata

COPY --from=build /app/product .
COPY --from=build /app/.env ./.env
COPY --from=build /app/public ./public

EXPOSE 8000

CMD [ "./product" ]
