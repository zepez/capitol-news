FROM golang:alpine

# important!
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS=-mod=vendor
ENV APP_USER app
ENV APP_HOME /go/src/microservices

# custom 
ENV cron * * * * *
ENV endpoint http://host.docker.internal:3001/
ENV seconds 60


RUN mkdir /app
ADD . /app
WORKDIR /app

# compile project
RUN go mod vendor
RUN go build

# open the port 8000
EXPOSE 3000
CMD [ "/app/capitol-news" ]










