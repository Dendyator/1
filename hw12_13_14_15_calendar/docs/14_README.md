сборка 

1) make build

2) make run

если подглючивает бд необходимо запустить из корня проекта команду
docker-compose -f deployments/docker-compose.yaml down --volumes --remove-orphans

3) после сборки применить миграции командой
docker exec -it calendar_service goose -dir /migrations postgres "postgres://user:password@db:5432/calendar?sslmode=disable" up

4) ввести события в терминале (команды представлены ниже)

5) make run-scheduler

6) make run-sender

RabbitMQ:
http://localhost:15672
guest/guest

время событий используется в формате UNIX, для удобства использования команд есть декодер времени
из обычного формата в формат UNIX hw12_13_14_15_calendar/cmd/unix-transformer/main.go

Использование gRPC
1. Создание события
grpcurl -plaintext -d '{
"event": {
"title": "Название события",
"description": "Описание события",
"startTime": 1731616200,
"endTime": 1731618000,
"userId": "e9b1f4b2-dc3e-4ea0-a8f3-1234567890ab"
}
}' localhost:50051 api.EventService/CreateEvent

2. Обновление события
   grpcurl -plaintext -d '{
   "id": "b1f4b2e9-dc3e-4ea0-a8f3-1234567890ab",
   "event": {
   "id": "b1f4b2e9-dc3e-4ea0-a8f3-1234567890ab",
   "title": "Обновлённое название",
   "description": "Обновлённое описание",
   "startTime": 1609459200,
   "endTime": 1609462800,
   "userId": "e9b1f4b2-dc3e-4ea0-a8f3-1234567890ab"
   }
   }' localhost:50051 api.EventService/UpdateEvent

3. Удаление события
   grpcurl -plaintext -d '{
   "id": "b1f4b2e9-dc3e-4ea0-a8f3-1234567890ab"
   }' localhost:50051 api.EventService/DeleteEvent

4. Получение события
   grpcurl -plaintext -d '{
   "id": "b1f4b2e9-dc3e-4ea0-a8f3-1234567890ab"
   }' localhost:50051 api.EventService/GetEvent

5. Список всех событий
   grpcurl -plaintext -d '{}' localhost:50051 api.EventService/ListEvents

6. Список событий за день
   grpcurl -plaintext -d '{
   "date": 1609459200  // Время в формате Unix для конкретного дня
   }' localhost:50051 api.EventService/ListEventsByDay

7. Список событий за неделю
   grpcurl -plaintext -d '{
   "start": 1609459200  // Время в формате Unix начала недели
   }' localhost:50051 api.EventService/ListEventsByWeek

8. Список событий за месяц // Время в формате Unix начала месяца
   grpcurl -plaintext -d '{
   "start": 1609459200  
   }' localhost:50051 api.EventService/ListEventsByMonth


