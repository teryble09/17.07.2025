# 17.07.2025
Тестовое задание, суть задания в tz.txt

Для запуска:
- go mod downoload
- go run cmd/archiver/main.go

Есть config.yaml в корне проекта, который включает в себя порт, разрешенные типы, время таймаута для скачивания, количество попыток скачивания и т.д.

Реализованы эндпоинты:
- post /tasks/ создать задачу, возвращает id задачи

Например:

{
  "task_id": "cf8af25b-493f-4d6e-aa85-80327a3ccaff"
}

- post /tasks/{id}/urls добавляет url для скачивания

Например:

{
  "url":"https://upload.wikimedia.org/wikipedia/commons/4/47/PNG_transparency_demonstration_1.png"
}

- get /tasks/{id}/status возвращает статус задачи

Например:

{
    "urls": [
        {
            "address": "https://upload.wikimedia.org/wikipedia/commons/4/47/PNG_transparency_demonstration_1.png",
            "status": "archived"
        }
    ]
}

- get /tasks/{id}/archive возвращает архив

Использованы паттерны:
 - reepository (см. internal/archiver/repository и internal/storage)
 - retry (см. internal/archiver/service/load.go)
 - semaphore (см. TaskService в internal/archiver/service/service.go)

Старался делать код как можно чище и делать код потокобезопасным с помощью использования mutex

telegram: https://t.me/teryble09
