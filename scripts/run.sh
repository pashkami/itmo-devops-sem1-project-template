set -e

echo "Компилируем приложение..."
go build -o app ./main.go

echo "Запускаем приложение..."
nohup ./app > output.log 2>&1 &

GO_PID=$!

echo "Приложение запущено с PID: $GO_PID"