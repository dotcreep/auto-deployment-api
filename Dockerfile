FROM golang:1.23-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN rm -rf archive config.yaml Makefile.bak
RUN go build -v -o /usr/local/bin/app ./...

VOLUME [ "storage/dist/" ]

CMD ["app"]