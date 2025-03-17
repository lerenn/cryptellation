![Cryptellation Logo](website/static/images/title.png)

Cryptellation is a **scalable cryptocurrency investment system**.

This system allows developers to create bots to manage their investments on
different cryptographic markets, featuring **backtesting**, **forward testing** and
**live running**.

## Serve the project with Docker

```bash
# Start the stack locally using Docker-compose
make docker/worker/up

# Check that everything works
go run ./cmd/cryptellation info
```

## Serve the project with Kind

```bash
# Start the stack locally using Kind
make kind/cryptellation/deploy

# Start port-forward
make kind/cryptellation/port-forward

# Check that everything works
go run ./cmd/cryptellation info
```