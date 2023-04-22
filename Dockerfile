FROM golang:latest

LABEL maintainer="Akmal Ergashev <ergashev2001@list.ru>"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 9000

RUN go build

RUN find . -name "*.go" -type f -delete

EXPOSE $PORT

CMD ["./SNS-connections"]