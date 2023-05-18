FROM golang:1.19
USER 10001
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /go-docker-demo
EXPOSE 8080
CMD [ "/go-docker-demo" ]