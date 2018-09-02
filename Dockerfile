FROM golang:latest
RUN apt-get update && apt-get install -y libmagickwand-dev

WORKDIR /go/src/github.com/Defman21/madnessBot
COPY . .
RUN pwd
RUN ls -la
ENV GO111MODULE=on
RUN go install -v ./...

EXPOSE 9000

CMD ["madnessBot"]
