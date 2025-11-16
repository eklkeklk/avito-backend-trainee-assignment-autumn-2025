# avito-backend-trainee-assignment-autumn-2025
test assignment for an internship position as a backend developer at Avito

## Инструкция по запуску
Перед запуском создайте файл '.env' на основе '.env.example':
```bash
cp .env.example .env
```

- Для запуска проекта необходимо перейти в корень проекта и ввести в консоли команду:
```bash
docker-compose up
```
Для ввода запросов необходимо открыть вторую консоль

- Для запуска проекта в фоновом режиме необходимо ввести в консоли команду:
```bash
docker-compose up -d
```
В таком случае, вводить запросы можно в той же консоли, в которой проект был запущен

#### Примеры запросов:
- Создание команды
```bash
curl -X 'POST' \
  'http://localhost:8080/team/add' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "team_name": "developers",
  "members": [
    {
      "user_id": "u3",
      "username": "Anna",
      "is_active": false
    },
    {
      "user_id": "u4",
      "username": "Beth",
      "is_active": true
    },
    {
      "user_id": "u5",
      "username": "Candice",
      "is_active": true
    },
    {
      "user_id": "u6",
      "username": "Danny",
      "is_active": true
    },
    {
      "user_id": "u7",
      "username": "Eugene",
      "is_active": false
    },
    {
      "user_id": "u8",
      "username": "Fogelle",
      "is_active": true
    },
    {
      "user_id": "u9",
      "username": "Gekkelle",
      "is_active": true
    }
  ]
}'
```

- Создание пулл-реквеста
```bash
curl -X 'POST' \
  'http://localhost:8080/pullRequest/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "pull_request_id": "pr-1002",
  "pull_request_name": "Fix",
  "author_id": "u4"
}'
```

- Переизбрание проверяющего
```bash
curl -X 'POST' \
  'http://localhost:8080/pullRequest/reassign' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "pull_request_id": "pr-1002",
  "old_reviewer_id": "u5"
}'
```

- Получить команду с участниками
```bash
curl -X 'GET' \
  'http://localhost:8080/team/get?team_name=developers' \
  -H 'accept: application/json'
```

- Получить пулл-реквесты, где пользователь назначен проверяющим
```bash
curl -X 'GET' \
  'http://localhost:8080/users/getReview?user_id=u6' \
  -H 'accept: application/json'
```

- Изменить статус пользователя
```bash
curl -X 'POST' \
  'http://localhost:8080/users/setIsActive' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "user_id": "u8",
  "is_active": false
}'
```

- Merge пулл-реквеста
```bash
curl -X 'POST' \
  'http://localhost:8080/pullRequest/merge' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "pull_request_id": "pr-1002"
}'
```

## Проблемы и решения
В условии не были описаны ошибки, возвращаемые при внутренних ошибках сервера, 
неверных запросах и неверных методах запроса.

Я добавила обработку этих ошибок, для того, чтобы пользователю было известно, по какой причине его запрос не прошел. Теперь в каждом из этих случаев выводится ответ:

- Внутренние ошибки сервера
```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "internal error"
  }
}
```
Код статуса: 500

- Ошибка при чтении тела запроса
```json
{
  "error": {
    "code": "BODY_READ_ERROR",
    "message": "body read error"
  }
}
```
Код статуса: 400

- Ошибка при парсинге json-реквеста
```json
{
  "error": {
    "code": "JSON_PARSE_ERROR",
    "message": "json unmarshall error"
  }
}
```
Код статуса: 400

- Неверный метод запроса
```json
{
  "error": {
    "code": "METHOD_NOT_ALLOWED",
    "message": "method not allowed"
  }
}
```
Код статуса: 405

- Неверный запрос (пустые значения необходимых полей)
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "invalid request"
  }
}
```
Код статуса: 400

## Дополнительные задания

### endpoint для статистики
Добавила простой endpoint статистики по просмотру числа назначений для пользователей и пулл-реквестов

#### Примеры запросов

- Статистика по пользователям
```bash
curl -X 'GET' \
  'http://localhost:8080/stats/users' \
  -H 'accept: application/json'
```

- Статистика по пулл-реквестам
```bash
curl -X 'GET' \
  'http://localhost:8080/stats/prs' \
  -H 'accept: application/json'
```

- Общая статистика
```bash
curl -X 'GET' \
  'http://localhost:8080/stats/overall' \
  -H 'accept: application/json'
```

### Описание линтера

Описание линтера можно найти в файле '.golangci.yml'

Локальный запуск линтера можно произвести командой:
```bash
make lint-conf
```

Или же

```bash
make all
```

Также есть стандартный линтер, запускаемый командой:
```bash
make lint
```