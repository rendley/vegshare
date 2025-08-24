FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Установка утилиты migrate
RUN apk --no-cache add curl &&     curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz &&     mv migrate /usr/local/bin/migrate &&     chmod +x /usr/local/bin/migrate

COPY . .

RUN go build -o /vegshare ./cmd/main.go

EXPOSE 8080

CMD [ "/vegshare" ]