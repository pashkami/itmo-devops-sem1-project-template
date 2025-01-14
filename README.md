# Финальный проект 1 семестра

## Поддерживаемые платформы
- **Операционные системы:** Windows, Linux, macOS
- **Go:** версия 1.20 и выше
- **PostgreSQL:** версия 13 и выше

---

## Требования к системе

### Аппаратные требования:
- **Оперативная память:** не менее 1 ГБ
- **Свободное дисковое пространство:** 250 МБ для базы данных

---

## Установка и запуск

1. Убедитесь, что установлены **Go** и **PostgreSQL**.
2. Склонируйте репозиторий:
   ```bash
   git clone <repository_url> && cd <repository_directory>
   ```
3. Создайте базу данных:
   ```bash
   psql -U validator -d postgres -c "CREATE DATABASE project_sem_1;"
   ```
4. Настройте базу данных, запустив скрипт:
   ```bash
   ./prepare.sh
   ```
5. Запустите сервер:
   ```bash
   ./run.sh
   ```
6. Сервер будет доступен по адресу: [http://localhost:8080](http://localhost:8080)

---

## Работа с данными

- **Пример данных:** Директория `sample_data` содержит разархивированную версию файла `sample_data.zip`.
- Убедитесь, что структура данных соответствует требованиям проекта.

---

## Тестирование

Для проверки работоспособности используйте скрипт `tests.sh`:
1. Отправка **POST-запроса** с файлом `sample_data.zip`.
2. Тестирование **GET-запроса** и проверка загрузки файла `data.zip`.

Запустите скрипт:
```bash
./tests.sh
```

---

## Контакты

Если возникнут вопросы, вы можете обратиться к **pashkami**.

