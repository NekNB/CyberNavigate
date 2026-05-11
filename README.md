# CyberNavigate

CyberNavigate - это круто

# Запуск

- `make run` - собирает весь проект и запускает его
- `make build SERVICE_NAME` - сборка образ
- `make stop SERVICE_NAME` - полная остановка сервиса
- `make update SERVICE_NAME` - обновление сервиса сервиса
- `make start` - полный запуск проекта из собранных ранее образов
- `make down` - полная остановка проекта
- `make scale SERVICE_NAME COUNT` - установка количества реплик для одного сервиса

- `make secrets_remove`|`make srm` - удаляет секреты
- `make secrets_create`|`make scr` - создает секреты
