FROM golang:1.8

RUN mkdir -p /cointron
WORKDIR /cointron
ADD . /cointron
RUN go get github.com/shopspring/decimal
RUN go get gopkg.in/jcelliott/turnpike.v2
RUN go get gopkg.in/telegram-bot-api.v4
CMD ["go", "run", "main.go"]