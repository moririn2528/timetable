FROM golang:1.18.1
RUN apt-get update && apt-get install git
RUN mkdir -p /api/go
WORKDIR /api/go
COPY api/go .
RUN mkdir -p /front
COPY front /front
RUN go get -d -v ./...
CMD ["go","run","main.go"]
