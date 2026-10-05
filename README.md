# go_projects

Учебный репозиторий с двумя сервисами на Go:

- **gateway** — простой HTTP-шлюз на порту `8080`.
- **ledger** — бизнес-логика: учёт финансовых транзакций с проверкой бюджетов по категориям.

У каждого сервиса собственный Go-модуль (`go.mod`), поэтому они собираются и запускаются независимо.

## Структура репозитория

    go_projects/
    ├── gateway/
    │   ├── go.mod
    │   └── main.go
    ├── ledger/
    │   ├── go.mod
    │   ├── main.go          # точка входа, демонстрация работы
    │   ├── transaction.go   # структура Transaction
    │   ├── storage.go       # хранилище транзакций и AddTransaction/ListTransactions
    │   ├── budget.go        # Budget, SetBudget, LoadBudgets
    │   └── budgets.json     # начальные бюджеты (читаются при запуске)
    └── README.md

## Требования

- Go 1.22 или выше

## Gateway

Простой HTTP-сервер. Пока реализован один маршрут.

### Запуск

    cd gateway
    go run .

Сервис слушает порт `8080`.

### Проверка

    curl -i http://localhost:8080/ping

Ожидаемый ответ:

    HTTP/1.1 200 OK
    Content-Type: text/plain; charset=utf-8
    Content-Length: 4

    pong

Остановить сервис — `Ctrl+C`.

## Ledger

Сервис для учёта транзакций и контроля бюджетов. Демонстрационная программа
в `main.go` показывает полный сценарий работы:

1. Загружает бюджеты из `budgets.json` через `LoadBudgets`.
2. Добавляет ещё один бюджет через `SetBudget` (категория `transport`).
3. Пытается добавить несколько транзакций:
   - обычные — успешно,
   - превышающую лимит — отклоняется с ошибкой `budget exceeded ...`,
   - с нулевой суммой — отклоняется с ошибкой `amount cannot be zero`.
4. Выводит список успешно сохранённых транзакций.

### Модель данных

```go
type Transaction struct {
    ID          int
    Amount      float64
    Category    string
    Description string
    Date        string
}

type Budget struct {
    Category string
    Limit    float64
    Period   string
}
```

### Основные функции

| Функция | Назначение |
|---|---|
| `AddTransaction(tx Transaction) error` | Добавляет транзакцию. Проверяет, что сумма положительна, и что суммарные траты по категории не превышают установленный бюджет. |
| `ListTransactions() []Transaction` | Возвращает копию списка сохранённых транзакций. |
| `SetBudget(b Budget) error` | Добавляет или обновляет бюджет для категории. |
| `GetBudget(category string) (Budget, bool)` | Возвращает бюджет по категории. |
| `LoadBudgets(r io.Reader) error` | Читает массив бюджетов из JSON и регистрирует их через `SetBudget`. |

### Формат `budgets.json`

    [
      {"category": "food", "limit": 5000, "period": "month"},
      {"category": "entertainment", "limit": 2000, "period": "month"}
    ]

### Запуск

    cd ledger
    go run .

### Ожидаемый вывод

    Ledger service started
    Budgets loaded from budgets.json
    Budget for transport set via SetBudget
    AddTransaction accepted: category=food amount=1500.50
    AddTransaction accepted: category=transport amount=300.00
    AddTransaction accepted: category=entertainment amount=75.20
    AddTransaction rejected (budget exceeded expected): budget exceeded for category "food": limit=5000.00, would be=5500.50
    AddTransaction rejected (zero amount expected): amount cannot be zero
    Transactions:
    ID=1, Amount=1500.50, Category=food, Description=Groceries, Date=2026-10-05
    ID=2, Amount=300.00, Category=transport, Description=Taxi, Date=2026-10-05
    ID=3, Amount=75.20, Category=entertainment, Description=Cinema, Date=2026-10-04

Обратите внимание: транзакция, превышающая бюджет, **не попадает** в итоговый список.

## Правила проверки бюджета

При добавлении транзакции `AddTransaction`:

1. Проверяет, что `Amount > 0`.
2. Ищет бюджет для `tx.Category`. Если бюджета нет — транзакция принимается без ограничений.
3. Если бюджет есть — суммирует все уже сохранённые транзакции этой категории
   и сравнивает с лимитом. Если `currentSum + tx.Amount > limit`, возвращается ошибка
   вида `budget exceeded for category "..."`, и транзакция не сохраняется.

## Сборка

Каждый сервис собирается отдельно, из своей папки:

    cd gateway && go build ./...
    cd ../ledger && go build ./...