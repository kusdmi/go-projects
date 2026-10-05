# go-projects

Репозиторий содержит два сервиса на Go:

- `gateway` — HTTP-шлюз.
- `ledger` — бизнес-логика для работы с транзакциями.

## Требования

- Go 1.22 или выше

## Запуск Gateway

    cd gateway
    go run .

Сервис слушает порт 8080. Проверка:

    curl -i http://localhost:8080/ping

Ожидаемый ответ: pong со статусом 200.

## Запуск Ledger

    cd ledger
    go run .

Ожидаемый вывод:

    Ledger service started
    AddTransaction error: amount cannot be zero
    Transactions:
    ID=1, Amount=1500.50, Category=food, Description=Groceries, Date=2026-10-05
    ID=2, Amount=300.00, Category=transport, Description=Taxi, Date=2026-10-05
    ID=3, Amount=75.20, Category=entertainment, Description=Cinema, Date=2026-10-04
