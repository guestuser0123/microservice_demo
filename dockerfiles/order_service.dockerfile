FROM golang:1.18
ADD ./order_service /go/src/order_service
WORKDIR /go/src/order_service
RUN go get order_service
RUN go install
EXPOSE ${ORDER_SERVICE_PORT}
CMD ["/go/bin/order_service"]
