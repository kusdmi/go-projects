# go_projects

Учебный репозиторий с двумя сервисами на Go:

- **gateway** — HTTP-шлюз, принимает запросы от клиентов.
- **ledger** — бизнес-логика: хранит данные в PostgreSQL, кэширует отчёты в Redis.

У каждого сервиса собственный Go-модуль (`go.mod`). Связь между ними настроена через `go.work` и локальный `replace` в `gateway/go.mod`.

---

## Требования

- Go 1.22+
- Docker (для запуска PostgreSQL и Redis)
- Утилита миграций [Goose](https://github.com/pressly/goose)

---

## 1. Запуск инфраструктуры

Запустите PostgreSQL (на порту `5433`) и Redis (на порту `6379`) через Docker:

```bash
docker run -d --name cashcraft-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=cashapp \
  -p 5433:5432 postgres:15

docker run -d --name cashcraft-redis \
  -p 6379:6379 redis:7-alpine
```
