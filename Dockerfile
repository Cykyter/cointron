
FROM golang:1.8

WORKDIR /go/src/cointron
COPY . .
RUN go-wrapper download
RUN go-wrapper install
CMD ["go-wrapper", "run"] # ["main"]