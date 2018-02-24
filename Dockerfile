FROM golang:latest

WORKDIR /go/src/github.com/Defman21/madnessBot
COPY . .
RUN pwd
RUN ls -la
RUN go install -v ./...

EXPOSE 9000

CMD ["madnessBot"]
