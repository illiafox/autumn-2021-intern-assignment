# Тестовое задание на позицию стажера-бекендера

---

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
Порт по умолчанию `8080`, Endpoint `http://localhost:8080/`

### `/get` - получение баланса пользователя
* Метод: `GET`
* Успешный Status code: `200 Accepted`
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

```json
{
   "user_id": 10, // Айди пользователя
   "base": "EUR"
}
```
Ответ:
```json
{
    "ok": true,
    "base": "EUR", // Валюта в которой представлен баланс
    "balance": "1.20" // 1.20 евро, НЕ в копейках
}
```
Возможная ошибка (включая из стандартного запроса):
```json
{
  "ok": false, // валюта не поддерживается
  "err": "base: abbreviation 'EUR' is not supported"
}
```
---
### `/change` - изменить баланс пользователя
* Метод: `GET`
* Успешный Status code: `200 Accepted`
#### 1. Стандартный запрос:

```json
{
   "user_id": 10, // Айди пользователя
   "change": 3000 // Сумма изменения баланса, В КОПЕЙКАХ
   "decription": "покупка хлеба" // описание транзакции
}
```