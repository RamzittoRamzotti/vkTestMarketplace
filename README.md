# vkTestMarketplace

Маркетплейс с REST API на Go

---

База данных: SQLite


---

## Запуск

1. **Пример запуска контейнера:**
   ```sh
   docker build -t vktestmarketplace .
   docker run --rm -p 8080:8080 vktestmarketplace
   ```
2. **Переменные окружения (.env):**
   ```env
   JWT_SECRET=your_secret_key
   DB_PATH=./storage/storage.db
   ```

---

## API

### Регистрация
- **POST /register**
- Тело:
  ```json
  {
    "login": "user1",
    "password": "password123"
  }
  ```
- Ответ: данные пользователя (без пароля)

### Авторизация
- **POST /login**
- Тело:
  ```json
  {
    "login": "user1",
    "password": "password123"
  }
  ```
- Ответ:
  ```json
  { "token": "...jwt..." }
  ```

### Размещение объявления
- **POST /ads**
- Заголовок: `Authorization: Bearer <token>`
- Тело:
  ```json
  {
    "title": "Книга",
    "description": "Почти новая",
    "image_url": "https://.../image.jpg",
    "price": 500
  }
  ```
- Ответ:
  ```json
  {
    "id": 1,
    "title": "Книга",
    "description": "Почти новая",
    "image_url": "https://.../image.jpg",
    "price": 500,
    "author": "user1",
    "is_mine": true
  }
  ```

### Лента объявлений
- **GET /ads**
- Параметры: `page`, `limit`, `sort_by` (`created_at`/`price`), `sort_order` (`asc`/`desc`), `min_price`, `max_price`
- Можно передать токен (опционально)
- Ответ:
  ```json
  [
    {
      "id": 1,
      "title": "Книга",
      "description": "Почти новая",
      "image_url": "https://.../image.jpg",
      "price": 500,
      "author": "user1",
      "is_mine": true
    },
    ...
  ]
  ```

---

## Валидация
- Логин: 3-32 символа, буквы/цифры/подчёркивания
- Пароль: 6-64 символов
- Заголовок: 3-100 символов
- Описание: 50-1000 символов
- Цена: >0 и <10 000 000
- Картинка: URL, заканчивающийся на .jpg/.jpeg/.png/.gif

---

