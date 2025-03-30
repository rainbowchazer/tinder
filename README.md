# 🖤 Tinder-Go: Бэкенд для дейтинг-приложения  

Этот проект — серверная часть простого дейтинг-приложения, написанная на **Go** с использованием **PostgreSQL**, **Gorilla Mux** и **WebSockets** для чатов.  

## Функционал  
✅ Регистрация и авторизация пользователей  
✅ Работа с профилем: просмотр, обновление, удаление  
✅ Система лайков и мэтчинга  
✅ Чат между совпавшими пользователями (WebSockets)  
✅ Контейнеризация через Docker  

---

## Запуск проекта  

### 1. Клонирование репозитория  
git clone https://github.com/твой-ник/tinder.git
cd tinder-go

2. Настройка базы данных
Вариант 1: через Docker (рекомендуется)
Просто запусти контейнеры:

docker-compose up -d
Это автоматически создаст PostgreSQL и применит init.sql для создания таблиц.

Вариант 2: вручную
Если у тебя уже есть PostgreSQL, создай базу и заполни её:

psql -U postgres -c "CREATE DATABASE tinder;"
psql -U postgres -d tinder -f db/init.sql

3. Запуск сервера

go run cmd/server.go
После этого сервер будет работать на http://localhost:8080.


API Эндпоинты

Регистрация
POST /register
Тело запроса (JSON):

{
  "email": "user@example.com",
  "password": "securepassword",
  "username": "JohnDoe",
  "age": 25
}

Логин
POST /login
Тело запроса (JSON):

{
  "email": "user@example.com",
  "password": "securepassword"
}
Ответ:

{
  "token": "jwt-токен"
}

Профиль
GET /profile — Получить профиль

PUT /profile/update — Обновить профиль

DELETE /profile/delete — Удалить профиль

Лайки и мэтчи
POST /like — Поставить лайк

GET /matches — Получить мэтчи

Чат по WebSocket
ws://localhost:8080/ws