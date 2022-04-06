# Тестовое задание на позицию стажера-бекендера

## Запуск

### 1. Требования
* **MySQL:** `8.0.28` и выше
* **Go:** `1.18`


### 2. Подготовка `MySQL`

Достаточно запустить файл `migrate.sql`
```sql
SOURCE migrate.sql
```
**docker-compose** сделает это автоматически

### 3. Сборка и запуск
```shell
go build -o server
./server
```
С нестандартным путем к конфиг файлу `config.toml`
```shell
./server -confing docker.toml
```
С нестандартным путем для лога `log.txt`
```shell
./server -log mylog.txt
```

## docker-compose
**API готова к работе** сразу после поднятия контейнера
### Запуск
```shell
docker-compose up
```
### Остановка
```shell
docker-compose down
```
Есть возможность дополнительной настройки использовав аргументы запуска.
Например, для подключения к локальной бд:
```yaml
environment:
  ARGS: -config config.toml -log mylog.txt
```

## Курсы валют

Для обновления валют был выбран сервис https://www.currencyconverterapi.com (рабочий ключ уже в конфиге)

Бесплатная версия имеет ограничение **150 запросов в секунду**, настройки вынесены в отдельную структуру:
```toml
[Exchanger]
Every=120 # Обновление каждые 2 минуты
Endpoint="https://free.currconv.com/api/v7/convert"
Key="b904b50e29877e8c38e0"
Bases=["EUR","USD"]
```
Для загрузки всех курсов:
```toml
Bases = [] # Или убрать поле
```
Для пропуска загрузки:
```toml
Skip = true
```
Ручное добавление/обновление через `exchange.Add`
```go
exchange.Add("EUR", 92.39)
```
## Структуры таблиц
### balances: 
```mysql
balance_id integer primary key -- айди баланса, можно изменить
user_id integer -- айди пользователя, можно изменить
balance integer  -- баланс В КОПЕЙКАХ
```
При изменении айди баланса также обновляется в транзакциях 
### transactions: 
```mysql
transaction_id integer primary key -- айди транзакции
balance_id integer -- айди баланса получателя
from_id integer -- айди баланса отправителя, не NULL только в переводах
action integer -- сумма изменения счета В КОПЕЙКАХ
date timestamp -- дата транзакции
description tinytext -- описание транзакции
```
`balance_id` и `from_id` являются внешними ключами, поэтому для удаления баланса была добавлена функция (и API к ней): `Delete(balanceID, userID int64) error`
```go
db.Delete(0, 10) // Удалить баланс с user_id 10
db.Delete(10, 0) // Удалить баланс с balance_id 10
db.Delete(10,10) // Удалить баланс с balance_id 10
```
Или использовать SQL запрос:
```sql
SET FOREIGN_KEY_CHECKS=0;

DELETE FROM balances WHERE balance_id = 10 -- через balance_id
DELETE FROM balances WHERE user_id = 10 -- через user_id

SET FOREIGN_KEY_CHECKS=1;
```

---
## Методы API
**Порт по умолчанию `8080`, Endpoint `http://localhost:8080/`**

Ошибки наполнения полей одинаковые, ниже будут рассмотрены только важные
### `/get` - получение баланса пользователя
* Метод: `GET`
* Успешный Status Code: `200 Accepted`
#### 1. Стандартный запрос:

```json
{
   "user_id": 10 // Айди пользователя
}
```
Ответ:
```json
{
    "ok": true,
    "base": "RUB", // Валюта в которой представлен баланс
    "balance": "111.00" // 111 рублей, НЕ в копейках
}
```
Возможная ошибка:
```json
{
  "ok": false, // баланс с таким user_id не найден
  "err": "get balance: balance with user id 10 not found"
}
```

#### 2. С конвертацией в другую алюту:

```json lines
{
   "user_id": 10, // Айди пользователя
   "base": "EUR"
}
```
Ответ:
```json lines
{
    "ok": true,
    "base": "EUR", // Валюта в которой представлен баланс
    "balance": "1.20" // 1.20 евро, НЕ в копейках
}
```
#### Возможная ошибка:
```json lines
{
  "ok": false, // валюта не поддерживается
  "err": "base: abbreviation 'EUR' is not supported"
}
```
---
### `/change` - изменить баланс пользователя
* Метод: `POST`
* Успешный Status Code: `200 Accepted`
#### 1. Стандартный запрос: Пополнение баланса

```json lines
{
   "user_id": 10, // Айди пользователя
   "change": 3000, // Сумма изменения баланса, В КОПЕЙКАХ
   "decription": "вернуть через день" // описание транзакции
}
```
Ответ:

**Если баланс с юзером не существует и `change > 0`, создается новый**

```json lines
{ "ok": true }
```
### 2. Снятие денег с баланса
```json lines
{
   "user_id": 10, // Айди пользователя
   "change": -3000, // Сумма снятия, В КОПЕЙКАХ
   "decription": "покупка хлеба" // описание транзакции
}
```
Ответ 
```json lines
{ "ok": true }
```
#### Возможная ошибка:
```json lines
{
  "ok": false, // если change < 0 и создается новый аккаунт
  "err": "change balance: change (-10000) is below zero, balance creating is forbidden"
}
```

---

### `/transfer` - перевод денег между балансами
* Метод: `POST`
* Успешный Status Code: `200 Accepted`
#### 1. Стандартный запрос:
```json lines
{
  "to_id": 20, // user_id получателя
  "from_id": 10, // user_id отправителя
  "amount": 300, // сумма перевода В КОПЕЙКАХ > 0
  "description": "вернул долг" // описание транзакции
}
```
**Если баланса с юзером `to_id` не существует, создается новый**

#### Возможные ошибки:
```json lines
{
  "ok": false, // from_id не хватает средств для перевода
  "err": "transfer: insufficient funds: missing 92.00"
}
```
---

### `/view` - просмотр транзакций пользователя
* Метод: `GET`
* Успешный Status Code: `200 Accepted`
#### 1. Стандартный запрос:
```json lines
{
  "user_id": 10, // Айди пользователя
  "sort": "DATE_DESC", // Тип сортировки
  "limit": 100 // Количество транзакций
}
```
**Поддерживаемые типы сортировки:**
```json lines
"DATE_DESC": От старых до новых
"DATE_ASC": От новых до старых
"SUM_DESC": От дорогих до дешевых
"SUM_ASC": От дешевых до дорогих
```
Ответ (вывод сокращен):
```json lines
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
      "transaction_id": 7, // Айди транзакции
      "balance_id": 1, // Айди баланса
      "from_id": "",  // Простое зачисление ИЛИ снятие
      "action": -10000, // Сумма изменения В КОПЕЙКАХ
      "date": "2022-04-05T19:46:32+03:00", // Время и дата транзакции
      "description": "долг" // Описание транзакции 
    }
  ]
}
```
Если транзакций не было, но баланс создан:
```json lines
{
  "ok": true,
  "transactions": null
}
```
### 2. Пагинация (сдвиг)
```json lines
{
  "user_id": 10, // Айди пользователя
  "sort": "SUM_ASC", // Тип сортировки
  "limit": 100, // Количество транзакций
  "offset": 2 // сдвиг вывода
}
```
Ответ идентичен с прошлым запросом

#### Возможные ошибки
```json lines
{
  "ok": false, // Баланс с user_id 10 не найден
  "err": "get transfers: get balance (id 10): balance with user id 10 not found"
}
```
### `/delete` - удаление баланса
Был добавлен из-за внешних ключей в таблице `transactions`
* Метод: `POST`
* Успешный Status Code: `200 Accepted`
#### Запросы:
```json lines
{
    "user_id":10 // Удалить баланс с user_id 10
}
```
```json lines
{
    "balance_id":10 // Удалить баланс с id 10
}
```
```json lines
{
    "balance_id":10, // Удалить баланс с id 10
    "user_id":10 // Игнорируется
}
```
#### Ответ:
```json lines
{
  "ok": true
}
```
#### Ошибки:
```json lines
{
  "ok": false, // Баланс с user_id 10 не найден
  "err": "delete balance (user_id 10, balance_id 0): get balance (userID 10): balance with user id 10 not found"
}
```

```json lines
{
  "ok": false, // Баланс с balance_id 10 не найден
  "err": "delete balance (user_id 0, balance_id 10): balance with ID 10 not found"
}
```

### Общие ошибки
* **decoding json**
* **проверка синтаксиса**
* **проверка айди (валидные > 0)**
* **encoding json**

## TODO:
1. Смена айди пользователя, иными словами передача баланса
2. Тесты SQL методов и самой API
