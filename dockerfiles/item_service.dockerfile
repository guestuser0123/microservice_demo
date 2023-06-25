FROM golang:1.18
ADD ./item_service /go/src/item_service
WORKDIR /go/src/item_service
RUN go get item_service
RUN go install
EXPOSE ${ITEM_SERVICE_PORT}
CMD ["/go/bin/item_service"]
