# CRUD Приложение для Управления Списком Книг
## Пример решения тестового задания для практического проекта #1 курса GOLANG NINJA

### Стэк
- go 1.17
- postgres

### Запуск
```docker-compose up --build```

## API

### POST /books
Create new book

##### Example Input:
```json
{
    "publish_date": "2022-04-05T16:32:20Z",
    "title": "test2",
    "author": "author2",
    "rating": 2
}
```

### GET /books
Get all books

##### Example Input:
```json
[
  {
    "id": 1,
    "publish_date": "2022-04-04T16:32:20Z",
    "title": "test2",
    "author": "author",
    "rating": 2
  },
  {
    "id": 2,
    "publish_date": "2022-04-05T16:32:20Z",
    "title": "test2",
    "author": "author2",
    "rating": 2
  }
]
```
### GET /book?id=1
Get book by ID

#### Example Output:
```json
{
    "id": 1,
    "publish_date": "2022-04-04T16:32:20Z",
    "title": "test2",
    "author": "author",
    "rating": 2
}
```

### UPDATE /book?id=1
Update book by ID

#### Example Input:
```json
{
  "publish_date": "2022-04-04T16:32:20Z",
  "title": "test1",
  "author": "author1",
  "rating": 2
}
```

### DELETE /book?id=1
Delete book by ID