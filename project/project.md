# Проектная работа

Тема: 
- Использование ClickHouse для исследования развития платежной системы "Оплати"
_____
Шаги выполнения:
- Создание датасетов в Clickhouse
- Генерация юзеров, сервисов(примерно 800сервисов по 5 категориями, около 1.5млн юзеров и от 5 до 200 ивентов на использование)
- Миграция из MySQL в Clickhouse
- Вычитывание с помощью Clickhouse событий из RabbitMQ
- Создание паблишера, который будет генерировать ивенты в RabbitMQ
- Создание дешбордов и графиков для аналитики данных с целью дальнейшего развития системы
____
Пометки:
- На момент генерации иветнов должна будет закладываться логика, которую можно было бы проследить посредством анализа данных с помощью Clickhouse


# Ход работы:

1. Разворачиваем все сервисы. При этом будет запущена первичная генерация ивентов и юзеров в контейнер mysql:
```bash
make up
```
2. Мигрируем базы данных из mysql в clickhouse. Предусматриваем дальнейший анализ ивентов и юзеров по времени. Соответственно на момент создания таблиц делаем партиционирование ивентов по месяцам и юзеров по возрасту.

```sql
CREATE TABLE events (
    id UInt64,
    price UInt32,
    service_id UInt32,
    location_id UInt32,
    timestamp DateTime
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, id)
SETTINGS index_granularity = 8192;

CREATE TABLE users (
    id UInt64,
    last_name String,
    first_name String,
    father_name Nullable(String),
    years_old UInt8,
    gender_code String,
    bank String
)
ENGINE = MergeTree()
PARTITION BY years_old
ORDER BY (years_old, id)
SETTINGS index_granularity = 8192;
```
Мигрируем event-ы и user-ов
```sql
INSERT INTO events SELECT
    id,
    price,
    service_id,
    location_id,
    timestamp
FROM mysql('mysql', 'appdb', 'events', 'root', 'qwerty')

Query id: 5b5b88f1-c994-434d-8523-7bf75f993f1e

Ok.

0 rows in set. Elapsed: 326.547 sec. Processed 110.53 million rows, 2.65 GB (338.47 thousand rows/s., 8.12 MB/s.)
Peak memory usage: 91.32 MiB.

INSERT INTO users SELECT
    id,
    last_name,
    first_name,
    father_name,
    years_old,
    gender_code,
    bank
FROM mysql('mysql', 'appdb', 'users', 'root', 'qwerty')

Query id: 0ea08cc7-4250-4262-ad25-e22453b5248e

Ok.
```
3. Загружаем json-объекты с наименованиями сервисов в clickhouse:

Для начала создаем соответствующую таблицу:
```sql
CREATE TABLE services (
    id UInt64,
    Name String,
    Category String
)
ENGINE = MergeTree()
ORDER BY id;
```
Далее запускаем скрипт, который записывает все значения `JSONEachRow` и генерирует для них `id rowNumberInAllBlocks() + (SELECT max(id) FROM services) AS id`

4. Теперь давайте заведем наш паблишер. Он будет генерировать сообщения в rabbitMQ. Мы же в свою очередь будем вычитывать с помощью Clickhouse эти сообщения и пополнять таблицу events.
...

# Смотрим на дешборд

- Почему это интересно? Что мы собственно хотим увидить?

1. На дешборде можно наблюдать за тем, как бы произошло распределение покупок по месяцам.

2. Так же на что тратит та или иная возврастная аудитория деньги с помощью сервиса "Оплати".

3. В каких чаще используется

4. Какими банками пользуются пользователи

Анализ этой и другой информации может оказаться полезным потребителю.