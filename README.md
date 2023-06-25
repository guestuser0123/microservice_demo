# microservice_demo

## Description

This project is a proof of concept for a simple ecommerce demo application that uses a distributed microservice architecture. 
This application consists of 3 main services - Order, Item, Voucher. This application uses the Saga design pattern for distributed transactions
and Snowflake ID to generate unique transaction IDs.

There are 3 exposed APIs -
1. GET /item/get?id={item_id}

> **Description:** Gets item info by item ID.

2. GET /voucher/get?id={voucher_id}

> **Description:** Gets voucher info by voucher ID.

3. POST /order/place-order

> **Description:** Places an order that takes in item ID, item qty, voucher ID and returns the order information including 
a unique transaction ID. If the voucher ID is 0, it is assumed that no voucher was applied
for the order. It is also assumed that you can only apply 1 voucher per order, hence there is not input for voucher qty.


## Getting started

To run this application, please make sure that you have the following installed on your local machine:
1. Docker
2. Go 1.18+ (optional, only for local development)
3. wrk (optional, only for stress testing)

## Run

Most of the setup has been done via the `docker-compose.yml` file. This also includes loading the databases with some mock data.
You can refer to `db_init/{service}_db_init.sql` for more information. To run the application, you can simply execute -

```
docker compose up --build --detach
```

Or you may run the following if you do not need to rebuild the application -
```
docker compose start --detach
```
<br /><br />
Once the application is up, you can run the following commands to test the application.

1. Get item by item id -
```
curl --location --request GET 'localhost:8080/item/get?id=1' \
--header 'Content-Type: application/json' \
--data '{
    "item_id": 1,
    "item_qty": 2,
    "voucher_id": 1
}'
```

2. Get voucher -
```
curl --location --request GET 'localhost:8082/voucher/get?id=1' \
--header 'Content-Type: application/json' \
--data '{
    "voucher_id": 1
}'
```

3. Place order -
```
curl --location 'localhost:8081/place-order' \
--header 'Content-Type: application/json' \
--data '{
    "item_id": 1,
    "item_qty": 2,
    "voucher_id": 1
}'
```
<br /><br />
To stop and teardown the application, simply execute the following command -

```
docker compose down
```

## Stress testing

To stress test the application, you can use the [wrk benchmark tool](https://github.com/wg/wrk). Sample stress test scripts
have been provided in the `scripts/stress` folder. Once you have installed wrk, you can run the following commands to stress
test the various service endpoints -

```
wrk -t12 -c100 -d20s -s ./scripts/stress/place_order_unique.lua "http://localhost:8081/place-order"
```

```
wrk -t8 -c200 -d20s -s  ./scripts/stress/get_voucher_by_id.lua http://localhost:8082 
```

pprof has been added to each of the microservice `main.go` code. This means that you can run the following pprof command whilst
running the stress test, so that you may analyse the performance of the application.

To profile the item service -
```
go tool pprof http://localhost:8080/debug/pprof/profile 
```

To profile the order service -
```
go tool pprof http://localhost:8081/debug/pprof/profile 
```

To profile the voucher service -
```
go tool pprof http://localhost:8082/debug/pprof/profile 
```

> Note: You may need to change the port numbers if you have modified the default ports set in the `.env` file.

## Credits

Saga - https://towardsdev.com/saga-pattern-with-golang-examples-18aa39c2cc12

Snowflake ID - https://github.com/bwmarrin/snowflake

