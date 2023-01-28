FROM golang:1.19
RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go install

RUN go build  -o /cb-worker
CMD ["/cb-worker"]
