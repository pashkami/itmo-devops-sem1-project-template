#!/bin/bash 

# Для прерывания скрипта в случае возникновения ошибок
set -e

# Список для сопоставления с env
REQUIRED_VARS=("POSTGRES_HOST" "POSTGRES_PORT" "POSTGRES_USER" "POSTGRES_PASSWORD" "POSTGRES_DB")

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo "Отсутствует файл .env."
    exit 1
fi

for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Переменная $var не задана"
        exit 1
    fi
done

go mod tidy
echo "Зависимости установлены"

if ! command -v psql &> /dev/null
then
    echo "Клиент PostgreSQL не установлен."
    exit 1
fi

echo "Подключение к PostgreSQL и Создание таблицы prices в базе данных $POSTGRES_DB..."

PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_DB -c "
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,           
    created_at DATE NOT NULL,        
    name VARCHAR(255) NOT NULL,      
    category VARCHAR(255) NOT NULL,  
    price DECIMAL(10, 2) NOT NULL    
);"

echo "БД подготовлена успешно"
