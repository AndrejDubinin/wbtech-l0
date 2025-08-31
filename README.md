# wbtech-l0: Order

## Запуск проекта

### 1. Настройка переменных окружения
Все переменные окружения находятся в файле: \
`build/dev/.env`

Сначала необходимо скопировать пример файла:
```bash
cp build/dev/.env.example build/dev/.env
```

### 2. Запуск контейнеров

Выполните команду:
```bash
make compose-up
```

### 3. Доступ к сервисам

Фронтенд будет доступен по адресу: \
[http://localhost:${FRONTEND_HOST_PORT}](http://localhost:8080)

Kafka UI доступен по адресу: \
[http://localhost:${KAFKA_UI_PORT}](http://localhost:8082)

В интерфейсе можно отправлять сообщения в Kafka.

***
## Управление контейнерами

Остановка контейнеров
```bash
make compose-down
```

Перезапуск контейнеров
```bash
make compose-rs
```
