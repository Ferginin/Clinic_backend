FROM golang:1.25-alpine

RUN mkdir "medlife_backend"

ADD . /medlife_backend/

WORKDIR /medlife_backend/cmd/app

RUN go build -o main .

CMD ["./main"]