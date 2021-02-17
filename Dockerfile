FROM golang:1.14
RUN mkdir /app
WORKDIR /app
ADD . .

RUN go build -o reporter  powerbi-live-reporting/cmd/reporting

CMD ["./reporter"]
