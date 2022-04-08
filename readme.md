# Тестовое задание на позицию стажера-бекендера

## Запуск

### 1. Требования
* **PostgreSQL:** `14.2` и выше
* **Go:** `1.18`


### 2. Подготовка `PostgreSQL`

Достаточно запустить файл `migrate-up.sql`
```sql
SOURCE migrate-up.sql
```
**docker-compose** сделает это автоматически

---

#### УДАЛЕНИЕ ТАБЛИЦ
```sql
SOURCE migrate-down.sql
```
### 3. Сборка и запуск
```shell
git clone https://github.com/illiafox/autumn-2021-intern-assignment
```

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

**Включение**: `config.toml`
```toml
[Exchange]
Skip = false # true - отключить
```
Бесплатная версия имеет ограничение **150 запросов в секунду**, настройки вынесены в отдельную структуру:
```toml
[Exchanger]
Every=120 # Обновление каждые 2 минуты
Endpoint="https://free.currconv.com/api/v7/convert"
Key="b904b50e29877e8c38e0"
Bases=["EUR","USD"]
```
Для загрузки всех курсов (130 запросов):
```toml
Bases = [] # Или убрать поле
```
Для пропуска загрузки (`true` по умолчанию):
```toml
Skip = true # false - не пропускать
```
Отключить интервальные обновления:
```shell
./server -skip
```
### Загрузка \ Сохранение в файл
**При удачном запуске курс сохраняется в `currencies.json`**

Пропустить вызов API и загрузиться с файла:
```shell
./server -load
```
Задать нестандартный путь `.json`:
```shell
./server -curr my_currencies.json
```



**Ручное добавление/обновление через `exchange.Add`**
```go
exchange.Add("EUR", 92.39)
```
## ~~Тесты~~: в разработке
**Запуск частично готовых в папке `tests`:**
```shell
 go test -fuzz=FuzzDatabase
```

## Структуры таблиц
### balances: 
```mysql
balance_id bigserial primary key -- айди баланса, можно изменить
user_id bigint -- айди пользователя, можно изменить
balance integer  -- баланс В КОПЕЙКАХ
```
### transactions: 
```postgresql
transaction_id bigserial primary key -- айди транзакции
balance_id bigint -- айди баланса получателя
from_id bigint -- айди баланса отправителя, не NULL только в переводах
action integer -- сумма изменения счета В КОПЕЙКАХ
date timestamp -- дата транзакции
description text -- описание транзакции
```
---
## Методы API
**Порт по умолчанию `8080`, Endpoint `http://localhost:8080/`**

Ошибки наполнения полей одинаковые, ниже будут рассмотрены только важные
### `/get` - получение баланса пользователя
* Метод: `GET`
* Успешный Status Code: `200 Accepted`
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

#### 2. С конвертацией в другую алюту:

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
* Успешный Status Code: `200 Accepted`
#### 1. Стандартный запрос: Пополнение баланса

```json5
{
   "user_id": 10, // Айди пользователя
   "change": 3000, // Сумма изменения баланса, В КОПЕЙКАХ
   "decription": "вернуть через день" // описание транзакции
}
```
Ответ:

**Если баланс с юзером не существует и `change > 0`, создается новый**

```json5
{ "ok": true }
```
### 2. Снятие денег с баланса
```json5
{
   "user_id": 10, // Айди пользователя
   "change": -3000, // Сумма снятия, В КОПЕЙКАХ
   "decription": "покупка хлеба" // описание транзакции
}
```
Ответ 
```json5
{ "ok": true }
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
* Успешный Status Code: `200 Accepted`
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
* Успешный Status Code: `200 Accepted`
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
* Успешный Status Code: `200 Accepted`
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
  "ok": false, // баланс с user_id 12 уже есть, для передачи можно удалить через /delete
  "err": "db.Switch(old 10 - new 12): balance with user_id 12 already exists"
}
```

### Общие ошибки
* **decoding json**
* **проверка синтаксиса**
* * **проверка айди (валидные > 0)**
* **ошибки бд**
* * **предвиденные ошибки**
* * **сбой работы**
* **encoding json**

## TODO:
1. Тесты SQL методов и самой API
