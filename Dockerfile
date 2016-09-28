FROM golang:1.6

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
ENV GOPATH $GOPATH:/go/src/app

RUN git clone --depth=1 --branch=master https://github.com/swagger-api/swagger-ui.git /go/src/swagger-ui

COPY . /go/src/app

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]

