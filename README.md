# Alokin - Zota integration

This project integrates Zota's API service into an example merchant application. Although you can use the server
as-is by just using its API and manually opening the Deposit page to complete the transaction, there is an extremely
bare-bones frontend application that can be used to simulate a normal user flow. You can find the frontend application
[here](https://github.com/federlizer/alokin-zota-integration-frontend).

This server has been created as a proof of concept, so there are a handful of things that aren't production ready. The
[Caveats](#caveats) section explores some of those missing features.

## Installation

In the root directory, there is a `Dockerfile` that can be used to build an image for `alokin`. Make sure that the
environment variables in the Dockerfile are setup by following the [Configuration](#configuration) section. To install
the docker image run the following command (inside the project's root directory):

```bash
$ docker build -t alokin .
```

where `alokin` can be replaced with any arbitrary name you'd like. Then, to run the application, you should run the
following command:

```bash
$ docker run -it --rm --name alokin-server -p8080:8080 alokin
```

The above command will create a container using the `alokin` image, start it, map port `8080` to the host machine and
finally, remove it when it exists. At this point, the server should be started and you can confirm if you're able to
reach the server by making a GET request to `localhost:8080/ping`. Next, you can refer to the [Usage](#usage) section.

## Configuration

To configure the application, you will need to set the following environment variables in the `Dockerfile`:

```Dockerfile
# The secret API key provided by Zota
ENV ZOTA_SECRET_KEY=11111111-1111-1111-1111-111111111111
# The endpoint ID provided by Zota
ENV ZOTA_ENDPOINT_ID=111111
# The merchant ID provided by Zota
ENV ZOTA_MERCHANT_ID=EXAMPLE-MERCHANT-ID
# The base URL provided by Zota
ENV ZOTA_BASE_URL=https://api.zotapay-sandbox.com
```

## Usage

The alokin webserver exposes a very simple API. In total, it has three endpoints that can be used:

| Method | Endpoint | Description                     |
|--------|----------|---------------------------------|
| `GET`  | `/ping`  | Ping the server. Test endpoint. |
| `GET`  | `/order` | Get all saved orders.           |
| `POST` | `/order` | Make a new order.               |


#### GET /ping

The ping endpoint is used to confirm the server is running and responding to commands. The server should respond with
a JSON object that contains a `message` field with a value `"pong"`.

#### GET /order

This endpoint returns all orders that have been created since the start of the application. Check out
[Caveats](#caveats) for more information about how orders are "persisted".

#### POST /order

Use this endpoint to make a new order. The endpoint expects a `Content-Type` header of `application/json` and a JSON
body that includes the following fields:

```json
{
    "description": "The description of the order (max 128 characters)",
    "amount": 13.37
}
```

Once the request is accepted, the application should return a response that redirects you to Zota's deposit page,
where you can perform the actual transaction. At the same time, the server will start a goroutine that will continuously
query Zota's API (`Order Status` every 10 seconds) for the status of the deposit. Once a final order status is received
or a maximum number of retries have been reached (20), the goroutine will be stopped and the `internal` order will be
updated (i.e. if you call `GET /order` you should see the `paymentStatus` field change from `PENDING` to `APPROVED` or
`FAILED`)

#### Example usage flow

1. Get all current orders `GET /order` (should be empty at startup)
2. Create a new order `POST /order` (you should get redirected 302 Found to Zota's payment page)
3. Query all current orders `GET /order` (should include the new order we just created)
3. Complete/Fail payment process (You'll get redirected to 404 NotFound URL - [Deposit redirectUrl caveat](#deposit-redirecturl))
4. Requery all orders `GET /order` (should show the same order, but with a different `paymentStatus` field)

Keep in mind that depending on the timing, the last step (step 4) might take up to 10 seconds to properly show the
updated `paymentStatus` field for the created order - this happens due to
[Order Status flow implementation caveat](#order-status-flow-implementations).

## Run tests

Currently, there have been implemented sample tests for the `zota` and `api` packages. To run them, you can run the
following commands:

```bash
$ go test ./zota
$ go test ./api

# Or run all available tests
$ go test ./...
```

## Structure of the application

The application is separated into three main packages: the `internal` package, which contains the merchant's
internal logic, the `api` package, which contains the logic that sets up the API webserver and the `zota` package,
which contains the Zota related structs and logic.

The idea behind this separation is twofold:
1. Make sure that the data sent or received by Zota is isolated, in case the API changes (separate internal models from zota's request/response models)
2. Make it easier to switch out components from the application (e.g. if the API webserver needs to change from `gin` to some other library/framework)

The application's entry point is in the `cmd/alokin/main.go` file.

## Caveats

There are a few caveats about using the application that need to be mentioned. There are a lot of things that needs to
be implemented for a proper, production-ready application, however since this application is mostly a time-sensitive
proof of concept, some of the production-ready features have been left out.

#### Persistence

The first caveat is that the application doesn't implement a real persistence system. Since this is a PoC/MVP, the
persistence only occurs in-memory - within the application process' lifetime. Any orders that are created while the
application is running will be able to be tracked and displayed. However, as soon as the application is restarted,
all previous orders will be forgotten and you'll start from scratch.

Ideally, the application would implement a database connection where it can properly persist the orders and users
created, but due to time constraints, this was omitted.

#### Users

Next, there are no real users either. Whenever the POST `/order` endpoint is called, an example user (Nikola Velichkov)
is created with hard-coded values. This feature was also omitted due to time constraints. In a real-life, production
scenario, the application should implement an authentication mechansim, so that each individual user will be able to
make orders on their own.

#### Order Status flow implementations

In the [`Order Status` documentation](https://doc.zota.com/deposit/1.0/?shell#order-status-request) it is highly
recommended to implement both a callback handler and the `Order Status Polling` strategies to confirm user's deposits.
However, the `callback` handler implementation requires the `alokin` server to be publicly accessible with a domain
name. This wasn't possible at the time of writing this application, so only the polling strategy has been implemented
for the `Order Status` flow.

#### Deposit redirectUrl

Zota's [Deposit request](https://doc.zota.com/deposit/1.0/?shell#deposit-request) requires a `redirectUrl` parameter
that will be used to redirect the user when the transaction has been completed (regardless of status). This parameter
is currently hard-coded to `https://federlizer.com/deposit-completed`. However, since the application is not hosted on
a publicly reachable server you can expect Zota's payment page to redirect you to a page that returns a 404.
