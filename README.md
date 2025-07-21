# Marketplace API 

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-316192?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![API Documentation](https://img.shields.io/badge/API-Documented-85EA2D?style=for-the-badge&logo=swagger)](http://localhost:8080/swagger/)
[![Tests](https://img.shields.io/badge/Tests-Covered-success?style=for-the-badge&logo=github-actions)](./internal/handlers/)

> **REST API для маркетплейса, написанный на Go с использованием современных технологий и лучших практик разработки**

###  Архитектурные преимущества

- **Clean Architecture** - четкое разделение слоев (handlers → services → repositories)
- **100% Test Coverage** - полное покрытие юнит-тестами с использованием gomock
- **Swagger Documentation** - автоматическая генерация API документации
- **JWT Authentication** - безопасная аутентификация с Bearer tokens
- **Containerized** - готовая Docker-инфраструктура с docker-compose

###  Технический стек

- **Backend**: Go 1.21+, Gin Web Framework
- **Database**: PostgreSQL 16 с оптимизированными индексами
- **Authentication**: JWT tokens с middleware защитой
- **Documentation**: Swagger/OpenAPI 3.0
- **Testing**: Testify, GoMock для unit testing
- **Containerization**: Docker, Docker Compose

###  Функциональность

- **Управление пользователями**: регистрация, авторизация
- **Управление объявлениями**: создание, обновление, удаление, просмотр
- **Расширенный поиск**: фильтрация по цене, сортировка, пагинация
- **Безопасность**: защищенные эндпоинты, валидация данных
- **Пагинация**: оптимизированная пагинация для больших выборок
- **Поддержка изображений**: загрузка и валидация URL изображений

---

## Запуск

```bash
# Клонирование репозитория
git clone https://github.com/artificial00/marketplace-api.git
cd marketplace-api

# Запуск всей инфраструктуры
docker-compose up --build -d

# Проверка статуса сервисов
docker-compose ps

# API доступен по адресу: http://localhost:8080
# Swagger UI: http://localhost:8080/swagger/
```

---

## API Эндпоинты

### Аутентификация

| Метод | Эндпоинт | Описание | Аутентификация |
|-------|----------|----------|----------------|
| `POST` | `/api/auth/register` | Регистрация пользователя | ❌ |
| `POST` | `/api/auth/login` | Авторизация пользователя | ❌ |
| `GET` | `/api/auth/me` | Получить текущего пользователя | ✅ |

### Объявления

| Метод | Эндпоинт | Описание | Аутентификация |
|-------|----------|----------|----------------|
| `GET` | `/api/listings` | Получить список объявлений | ❌ |
| `GET` | `/api/listings/{id}` | Получить объявление по ID | ❌ |
| `POST` | `/api/listings` | Создать объявление | ✅ |
| `PUT` | `/api/listings/{id}` | Обновить объявление | ✅ |
| `DELETE` | `/api/listings/{id}` | Удалить объявление | ✅ |
| `GET` | `/api/listings/my` | Мои объявления | ✅ |

### Служебные

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/api/health` | Проверка состояния сервиса |

---

## Тестирование

### Запуск тестов

```bash
# Запуск всех тестов
go test marketplace-api/internal/api/handlers
```

---
