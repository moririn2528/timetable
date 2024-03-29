FROM golang:1.20
RUN apt-get update && apt-get install git
RUN mkdir -p /api/go
WORKDIR /api
RUN mkdir -p /front
COPY ./api/go/start.sh .
RUN chmod +x start.sh
COPY ./api/go /api/go
COPY ./front /front
ENV EXEC_ENV=docker
CMD ["./start.sh"]