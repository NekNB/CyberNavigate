# CyberNavigate
CyberNavigate - это круто



# Запуск
- `make run` - собирает весь проект и запускает его
- `make build` - сборка образов
- `make start` - полный запуск проекта
- `make stop` - полная остановка проекта


Swagger
 go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
 protoc -I . -I ./third_party --openapiv2_out=./swagger proto/user/user.proto