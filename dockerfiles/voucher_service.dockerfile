FROM golang:1.18
ADD ./voucher_service /go/src/voucher_service
WORKDIR /go/src/voucher_service
RUN go get voucher_service
RUN go install
EXPOSE ${VOUCHER_SERVICE_PORT}
CMD ["/go/bin/voucher_service"]
