## Домашнее задание №14 «Кроликизация Календаря»
Необходимо реализовать "напоминания" о событиях с помощью RabbitMQ (кролика).
Общая концепция описана в [техническом задании](./CALENDAR.MD).

Порядок выполнения ДЗ:
* установить локально очередь сообщений RabbitMQ (или сразу через Docker, если знаете как);
* создать процесс Планировщик (`scheduler`), который периодически сканирует основную базу данных,
выбирая события о которых нужно напомнить:
    - при запуске процесс должен подключаться к RabbitMQ и создавать все необходимые структуры
    (топики и пр.) в ней;
    - процесс должен выбирать сообытия для которых следует отправить уведомление (у события есть соотв. поле),
    создавать для каждого Уведомление (описание сущности см. в [ТЗ](./CALENDAR.MD)),
    сериализовать его (например, в JSON) и складывать в очередь;
    - процесс должен очищать старые (произошедшие более 1 года назад) события.
* создать процесс Рассыльщик (`sender`), который читает сообщения из очереди и шлёт уведомления;
непосредственно отправку делать не нужно - достаточно логировать сообщения / выводить в STDOUT.
* настройки подключения к очереди, периодичность запуска и пр. настройки процессов вынести в конфиг проекта;
* работу с кроликом вынести в отдельный пакет, который будут использовать пакеты, реализующие процессы выше.

Процессы не должны зависеть от конкретной реализации RMQ-клиента.

В результате компиляции проекта (`make build`) должно получаться 3 отдельных исполняемых файла
(по одному на микросервис):
- API (`calendar`);
- Планировщик (`calendar_scheduler`);
- Рассыльщик (`calendar_sender`).

Каждый из сервисов должен принимать путь файлу конфигурации:
```bash
./calendar           --config=/path/to/calendar_config.yaml
./calendar_scheduler --config=/path/to/scheduler_config.yaml
./calendar_sender    --config=/path/to/sender_config.yaml
```

После запуска RabbitMQ и PostgreSQL процессы `calendar_scheduler` и `calendar_sender`
должны запускаться без дополнительных действий.

### Критерии оценки
- Makefile заполнен и пайплайн зеленый - 1 балл
- Работа с RMQ выделена в отдельный пакет, код не дублируется - 1 балл
- Реализован Планировщик:
    - отсылает уведомления о выбранных событиях - 2 балла
    - удаляет старые события - 1 балл
- Реализован Рассыльщик - 2 балла
- Можно собрать сервисы одной командой (`make build`) - 1 балл
- Понятность и чистота кода - до 2 баллов

#### Зачёт от 7 баллов

За


Использование REST API с cURL

1) Добавить событие:

$event = @{
"id" = "12344";
"title" = "New Event 6";
"description" = "Event Description";
"start_time" = (Get-Date).AddHours(1).ToString("u");
"end_time" = (Get-Date).AddHours(2).ToString("u")
}

Invoke-RestMethod -Uri http://localhost:8080/events -Method POST -Body ($event | ConvertTo-Json) -ContentType "application/json"


2) Обновить событие:
$event = @{
"id" = "12345";
"title" = "Updated Event";
"description" = "Updated Description";
"start_time" = (Get-Date).AddHours(1).ToString("u");
"end_time" = (Get-Date).AddHours(3).ToString("u")
}

Invoke-RestMethod -Uri http://localhost:8080/events/12345 -Method PUT -Body ($event | ConvertTo-Json) -ContentType "application/json"

3) Получить событие:
Invoke-RestMethod -Uri http://localhost:8080/events/12345 -Method GET

4) Удалить событие:
Invoke-RestMethod -Uri http://localhost:8080/events/12345 -Method DELETE

5) Получить список событий:
   Invoke-RestMethod -Uri http://localhost:8080/events -Method GET



Использование gRPC

grpcurl -plaintext -d '{"event": {"id": "1", "title": "Test Event", "description": "This is a test", "start_time": 1700000000, "end_time": 1700007200}}' localhost:50051 api.EventService/CreateEvent


RabbitMQ:
http://localhost:15672
guest/guest