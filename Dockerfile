FROM golang:alpine
RUN apk update
RUN apk add imagemagick imagemagick-dev
RUN apk add git openssh gcc musl-dev

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go get

COPY . .

RUN go build -o /app/madnessBot

EXPOSE 9000
CMD ["/app/madnessBot", "-graphite"]
