Запуск: docker-compose up --build

frontend: http://localhost:3000
backend: http://localhost:8080

Админ:
admin
123

REST API
Auth:
POST /auth/register { "email": "example@mail.com", "password": "123456" }
POST /auth/login { "email": "example@mail.com", "password": "123456" }

Task:
admin/
POST /admin/task { "title": "Задача", "description": "Описание", "status": "new" }
PATCH /admin/assign { "taskID": 1, "userID": 2 }
DELETE /admin/task/{id} -
GET /admin/users -
Для всех:
GET /tasks -
PATCH /task/{id} { "status": "in_progress/done" }
