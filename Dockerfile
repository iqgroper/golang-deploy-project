FROM golang:alpine

COPY . /myapp

WORKDIR /myapp

RUN apk add make && apk add git && make build

CMD ["/myapp/bin/cruddapp"]