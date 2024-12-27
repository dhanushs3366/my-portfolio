# Dockerfile

FROM golang:1.23


WORKDIR /app/

COPY go.* .

RUN go  install github.com/air-verse/air@latest
RUN go mod download

COPY . .

EXPOSE 8080


CMD ["air","-c",".air.toml"]
