FROM golang:1.19
RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go install

RUN go run worker.go
#CMD ["/app/main"]
