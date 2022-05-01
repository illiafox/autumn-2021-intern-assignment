# Test task for the position of trainee golang backend developer

### [Версия на Русском](https://github.com/illiafox/autumn-2021-intern-assignment/blob/master/russian.md)

---
## Building / Running

### 1. Requirements

* **PostgreSQL:** `14.2`
* **Go:** `1.17`

### 2. Setuping `PostgreSQL`

#### Tables creation

*docker-compose does this automatically*

```shell
migrate -database ${POSTGRESQL_URL} -path migrate/ up
```

#### Dropping

```shell
migrate -database ${POSTGRESQL_URL} -path migrate/ down
```

### 3. Application building

```shell
git clone https://github.com/illiafox/autumn-2021-intern-assignment avito
cd avito/cmd/app

go build -o server
./server
```

#### With non-standard config and log file paths

```shell
server -confing config.toml -log log.txt
```

#### With reading from `environment`:

Available keys can be found in **[config structure](https://github.com/illiafox/autumn-2021-intern-assignment/blob/master/utils/config/struct.go#L3:L25)** tags

```shell
POSTGRES_PORT=4585 server -env
```

## docker-compose

**API starts** immediately after containers are up

```shell
docker-compose up
```

### Stopping

```shell
docker-compose down
```

---

It is possible to additionally configure the app using environment variables
```yaml
environment:
  POSTGRES_IP: 127.0.0.1 # connect to local database
  EXCHANGER_SKIP: true # skip currency loading
```

## Logs
In addition to the terminal output, logs are also written to the file
```shell
# Terminal
01/05/2022 14:54:38     info    Initializing database
```

```json5
// File (default log.txt)
{"level":"info","ts":"Sun, 01 May 2022 14:54:38 EEST","msg":"Initializing database"}
```



## Exchange Currencies

**Service:** `https://www.currencyconverterapi.com` (working key in the config)

**Enabling**: `config.toml`

```toml
[Exchange]
Skip = false # true - disable
```

The free version has a limitation of **150 requests per hour**,the settings are placed in the structure:

```toml
[Exchanger]
Every = 120 # Update every 2 minutes
Endpoint = "https://free.currconv.com/api/v7/convert"
Key = "b904b50e29877e8c38e0"
Bases = ["EUR", "USD"]
```

#### Load all available currencies (130 requests):

```toml
Bases = [] # Or remove field
```

### Loading \ Writing to file

#### If the loading is successfull, the rates are saved in `currencies.json`

---
#### Skip api loading

```toml
[Exchanger]
Skip = false
Load = true
```

#### Set saving path:

```toml
[Exchanger]
Skip = false
Path = "saves/currencies.toml"
```

**You can add currencies manually via `exchange.Add`**

```go
exchange.Add("EUR", 92.39)
```

## Tables structures

### balances:

```mysql
`balance_id` bigserial primary key -- balance id
`user_id bigint` -- user id, can be changed
`balance integer`  -- balance in cents
```

### transactions:

```mysql
`transaction_id` bigserial primary key -- transaction id
`balance_id` bigint -- receiver balance id
`from_id` bigint -- sender balance id, not NULL only in transfer operations
`action` integer -- change in cents
`date` timestamp -- transaction time
`description` text -- transaction description
    
 FOREIGN KEY (balance_id, from_id) REFERENCES balances (balance_id)
```

---

## API Methods

### ~~Tests~~: in development

**Default port `8080`, Endpoint `http://localhost:8080/`**

---

#### `200` Accepted - request was proceeded successfully

#### `400` Bad Request - wrong input data format

#### `406` Not Acceptable - data was not found

#### `422` Unprocessable Entity - request cannot be proceeded with current input data

#### `502` Internal Server Error - CRITICAL DATABASE ERROR

---

### `/get` - Get User Balance

* Method: `GET`

#### 1. Default request:

```json5
{
  "user_id": 10 // user id
}
```

Response:

```json5
{
  "ok": true,
  "base": "RUB", // Currency
  "balance": "111.00" // Balance NOT IN CENTS
}
```

Possible errors:

```json5
{
  "ok": false, // balance not found
  "err": "get balance: balance with user id 10 not found"
}
```

#### 2. Convert to another currency:

```json5
{
  "user_id": 10, // user id
  "base": "EUR" // currency abbreviation
}
```

Response:

```json5
{
  "ok": true,
  "base": "EUR", // Currency abbreviation
  "rate": 80.12, // Currency rate
  "balance": "1.20" // Balance NOT IN cents
}
```

#### Possible error:

```json5
{
  "ok": false, // Currency not available/supported
  "err": "base: abbreviation 'EUR' is not supported"
}
```

---

### `/change` - deposit/withdraw user balance

* Method: `POST`

#### 1. Default request: Deposit

```json5
{
  "user_id": 10, // user id
  "change": 3000, // change amount IN CENTS
  "decription": "salary" // transaction description
}
```

Response:

**If balance does not exists and `change > 0`,  new balance will be created**

```json5
{
  "ok": true
}
```

### 2. Withdraw 

```json5
{
  "user_id": 10, // user id
  "change": -3000, // amount IN CENTS
  "decription": "supermarket" // описание транзакции
}
```

Ответ

```json5
{
  "ok": true
}
```

#### Possible errors:

```json5
{
  "ok": false, // if change < 0 and balance not exist
  "err": "change balance: change (-10000) is below zero, balance creating is forbidden"
}
```

---

### `/transfer` - transfer money between balances

* Method: `POST`

#### 1. Default request:

```json5
{
  "to_id": 20, // receiver user id
  "from_id": 10, // sender user id
  "amount": 300, // amount IN CENTS > 0
  "description": "from Mark" // transaction description
}
```

**If balance `to_id` does not exist, new balance will be created**

#### Possible errors:

```json5
{
  "ok": false, // from_id has not enough money
  "err": "transfer: insufficient funds: missing 92.00"
}
```

---

### `/view` - view user transaction

* Method: `GET`

#### 1. Default request:

```json5
{
  "user_id": 10, // user id
  "sort": "DATE_DESC", // sorting type
  "limit": 100 // output size
}
```

**Sorting types supported:**

```json5
"DATE_DESC": From older to newer
"DATE_ASC": From newer to older
"SUM_DESC": From biggest transactions to lowest
"SUM_ASC": From lowest transactions to biggest
```

Response (output abbreviated):

```json5
{
  "ok": true,
  "transactions": [
    {
      "transaction_id": 8, // transaction id
      "balance_id": 1, // balance id
      "from_id": "12", // sender balance id
      "action": 2000, // amount IN CENTS
      "date": "2022-04-06T09:31:05+03:00", // transaction date
      "description": "долг" // transaction description
    },
    {
      "transaction_id": 7, 
      "balance_id": 1, 
      "from_id": "", // this is not transfer
      "action": -10000, 
      "date": "2022-04-05T19:46:32+03:00", 
      "description": "долг" 
    }
  ]
}
```

If there are no transactions, but balance exists:

```json5
{
  "ok": true,
  "transactions": null
}
```

### 2. Pagination (offset)

```json5
{
  "user_id": 10, // user id
  "sort": "SUM_ASC", // sorting type
  "limit": 100, // output size
  "offset": 2 // output offset
}
```

The response is the same as the in previous request, but with offset in 2 transactions

#### Possible errors

```json5
{
  "ok": false, // balance not found
  "err": "get transfers: get balance (id 10): balance with user id 10 not found"
}
```

### `/switch` - change balance user id

* Method: `POST`

#### Request:

```json5
{
  "old_user_id": 10, // old user id
  "new_user_id": 12 // new user id
}
```

#### Possible Errors:

If balance with new user id already exists
```json5
{
  "ok": false,
  "err": "db.Switch(old 10 - new 12): balance with user_id 12 already exists"
}
```

---

## Old branches:
* ### [FastHTTP router](https://github.com/illiafox/autumn-2021-intern-assignment/tree/fasthttp)
* ### [MySQL version](https://github.com/illiafox/autumn-2021-intern-assignment/tree/mysql)
