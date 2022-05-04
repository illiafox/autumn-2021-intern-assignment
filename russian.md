# Тестовое задание на позицию стажера-бекендера

### [English Translation](https://github.com/illiafox/autumn-2021-intern-assignment/blob/master/readme.md)

---

## Запуск

### 1. Требования

* **PostgreSQL:** `14.2`
* **Go:** `1.17`

### 2. Подготовка `PostgreSQL`

#### Создание таблиц

docker-compose сделает это автоматически

```shell
migrate -database ${POSTGRESQL_URL} -path migrate/ up
```

#### Удаление

```shell
migrate -database ${POSTGRESQL_URL} -path migrate/ down
```

### 3. Сборка и запуск

```shell
git clone https://github.com/illiafox/autumn-2021-intern-assignment avito
cd avito/cmd/app

go build -o server
./server
```

#### С нестандартными путями к конфигу и лог файлу

```shell
server -confing config.toml -log log.txt
```

#### С чтением `environment` переменных:

Доступные значения можно посмотреть в
**[тегах структуры](https://github.com/illiafox/autumn-2021-intern-assignment/blob/master/utils/config/erruct.go#L3:L25)** конфига

```shell
POSTGRES_PORT=4585 server -env
```

## docker-compose

**API готова к работе** после поднятия контейнеров

### Запуск

```shell
docker-compose up
```

### Остановка

```shell
docker-compose down
```

Есть возможность дополнительной настройки, используя аргументы запуска

```yaml
environment:
  POSTGRES_IP: 127.0.0.1 # подключение к локальной бд
  EXCHANGER_SKIP: true # пропуск загрузки валют 
```

## Логи
Кроме вывода в терминал, логи также пишутся в файл
```shell
# Терминал
01/05/2022 14:54:38     info    Initializing database
```

```json5
// Файл (log.txt по-умолчанию)
{"level":"info","ts":"Sun, 01 May 2022 14:54:38 EEST","msg":"Initializing database"}
```

---

## Тесты базы данных

### 1. Создайте новую бд и таблицы через миграцию

### 2. Запустите:

```shell
POSTGRES_DATABASE=avito_test go test
```

#### Тесты запустятся при условии, что таблица пуста

#### Таблицы очистятся автоматически

--- 
### Тесты API: скоро

---



## Курсы валют

Для обновления валют был выбран сервис https://www.currencyconverterapi.com (рабочий ключ уже в конфиге)

**Включение**: `config.toml`

```toml
[Exchange]
Skip = false # true - отключить
```

Бесплатная версия имеет ограничение **150 запросов в час**, настройки вынесены в отдельную структуру:

```toml
[Exchanger]
Every = 120 # Обновление каждые 2 минуты
Endpoint = "https://free.currconv.com/api/v7/convert"
Key = "b904b50e29877e8c38e0"
Bases = ["EUR", "USD"]
```


Для загрузки всех курсов (130 запросов):

```toml
Bases = [] # Или убрать поле
```

### Загрузка \ Сохранение в файл

#### При удачном запуске курс сохраняется в `currencies.json`

---

#### Пропустить вызов API и загрузиться с файла:

```toml
[Exchanger]
Skip = false
Load = true
```

#### Задать нестандартный путь `.json`:

```toml
[Exchanger]
Skip = false
Path = "saves/currencies.toml"
```

**Ручное добавление/обновление через `exchange.Add`**

```go
exchange.Add("EUR", 92.39)
```




## Структуры таблиц

### balances:

```mysql
`balance_id` bigserial primary key -- айди баланса
`user_id bigint` -- айди пользователя, можно изменить
`balance integer`  -- баланс В КОПЕЙКАХ
```

### transactions:

```mysql
`transaction_id` bigserial primary key -- айди транзакции
`balance_id` bigint -- айди баланса получателя
`from_id` bigint -- айди баланса отправителя, не NULL только в переводах
`action` integer -- сумма изменения счета В КОПЕЙКАХ
`date` timestamp -- дата транзакции
`description` text -- описание транзакции
```

---

## Методы API

**Порт по умолчанию `8080`, Endpoint `http://localhost:8080/`**

---

#### `200` Accepted - запрос успешно выполнен

#### `400` Bad Request - неверный формат входных данных

#### `406` Not Acceptable - данные не найдены

#### `422` Unprocessable Entity - запрос не может выполняться дальше с текущими входными значениями

#### `502` Internal Server Error - критическая ошибка (базы данных)

---

### `/get` - получение баланса пользователя

* Метод: `GET`

#### 1. Стандартный запрос:

```json5
{
  "user_id": 10 // Айди пользователя
}
```

Ответ:

```json5
{
  "ok": true,
  "base": "RUB", // Валюта в которой представлен баланс
  "balance": "111.00" // 111 рублей, НЕ в копейках
}
```

Возможная ошибка:

```json5
{
  "ok": false, // баланс с таким user_id не найден
  "err": "get balance: balance with user id 10 not found"
}
```

#### 2. С конвертацией в другую валюту:

```json5
{
  "user_id": 10, // Айди пользователя
  "base": "EUR"
}
```

Ответ:

```json5
{
  "ok": true,
  "base": "EUR", // Валюта в которой представлен баланс
  "rate": 80.12, // Курс валюты
  "balance": "1.20" // 1.20 евро, НЕ в копейках
}
```

#### Возможная ошибка:

```json5
{
  "ok": false, // валюта не поддерживается
  "err": "base: abbreviation 'EUR' is not supported"
}
```

---

### `/change` - изменить баланс пользователя

* Метод: `POST`

#### 1. Стандартный запрос: Пополнение баланса

```json5
{
  "user_id": 10, // Айди пользователя
  "change": 3000, // Сумма изменения баланса, В КОПЕЙКАХ
  "description": "вернуть через день" // описание транзакции
}
```

Ответ:

**Если баланс с юзером не существует и `change > 0`, создается новый**

```json5
{
  "ok": true
}
```

### 2. Снятие денег с баланса

```json5
{
  "user_id": 10, // Айди пользователя
  "change": -3000, // Сумма снятия, В КОПЕЙКАХ
  "description": "покупка хлеба" // описание транзакции
}
```

Ответ

```json5
{
  "ok": true
}
```

#### Возможная ошибка:

```json5
{
  "ok": false, // если change < 0 и создается новый аккаунт
  "err": "change balance: change (-10000) is below zero, balance creating is forbidden"
}
```

---

### `/transfer` - перевод денег между балансами

* Метод: `POST`

#### 1. Стандартный запрос:

```json5
{
  "to_id": 20, // user_id получателя
  "from_id": 10, // user_id отправителя
  "amount": 300, // сумма перевода В КОПЕЙКАХ > 0
  "description": "вернул долг" // описание транзакции
}
```

**Если баланса с юзером `to_id` не существует, создается новый**

#### Возможные ошибки:

```json5
{
  "ok": false, // from_id не хватает средств для перевода
  "err": "transfer: insufficient funds: missing 92.00"
}
```

---

### `/view` - просмотр транзакций пользователя

* Метод: `GET`

#### 1. Стандартный запрос:

```json5
{
  "user_id": 10, // Айди пользователя
  "sort": "DATE_DESC", // Тип сортировки
  "limit": 100 // Количество транзакций
}
```

**Поддерживаемые типы сортировки:**

```json5
"DATE_DESC": От старых транзакций до новых
"DATE_ASC": От новых до старых
"SUM_DESC": От больших сделок до маленьких
"SUM_ASC": От маленьких до больших
```

Ответ (вывод сокращен):

```json5
{
  "ok": true,
  "transactions": [
    {
      "transaction_id": 8, // Айди транзакции
      "balance_id": 1, // Айди баланса
      "from_id": "12", // Перевод получен от balance_id 12
      "action": 2000, // Сумма перевода В КОПЕЙКАХ
      "date": "2022-04-06T09:31:05+03:00", // Время и дата транзакции
      "description": "долг" // Описание транзакции 
    },
    {
      "transaction_id": 7,
      "balance_id": 1, 
      "from_id": "", // Простое зачисление ИЛИ снятие
      "action": -10000, 
      "date": "2022-04-05T19:46:32+03:00", 
      "description": "долг" 
      
    }
  ]
}
```

Если транзакций не было, но баланс создан:

```json5
{
  "ok": true,
  "transactions": null
}
```

### 2. Пагинация (сдвиг)

```json5
{
  "user_id": 10, // Айди пользователя
  "sort": "SUM_ASC", // Тип сортировки
  "limit": 100, // Количество транзакций
  "offset": 2 // сдвиг вывода
}
```

Ответ идентичен с прошлым запросом, но с выполнением сдвига на 2 транзакции

#### Возможные ошибки

```json5
{
  "ok": false, // Баланс с user_id 10 не найден
  "err": "get transfers: get balance (id 10): balance with user id 10 not found"
}
```

### `/switch` - передача баланса другому пользователю

* Метод: `POST`

#### Запрос:

```json5
{
  "old_user_id": 10, // id прошлого владельца
  "new_user_id": 12 // id нового владельца
}
```

#### Ошибки:

Если баланс нового владельца уже существует

```json5
{
  "ok": false, // баланс с user_id 12 уже есть, для передачи можно удалить через /del
  "err": "tx.Switch(old 10 - new 12): balance with user_id 12 already exists"
}
```

---

### `/delete` - удалить баланс

#### Транзакции не будут удалены

* Метод: `POST`

#### Запрос::

```json5
{
  "user_id": 10, // user id
}
```

#### Возможная ошибка:

```json5
{
  "ok": false, // баланс не найдет
  "err": "delete balance: balance with user id 10 not found"
}
```



---



## Старые версии:
* ### [FastHTTP router](https://github.com/illiafox/autumn-2021-intern-assignment/tree/fasthttp)
* ### [Версия для MySQL](https://github.com/illiafox/autumn-2021-intern-assignment/tree/mysql)
