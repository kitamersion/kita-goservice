# Kita Go Service

Microservice for GO

## Getting Started

If you use nix: `nix-shell` or `nix-shell shell.nix`

**Start infrastructure**

```bash
make docker-up
```

**Create Kafka topics**

```bash
make create-topics
```

**Run applications**

```bash
# Terminal One - Public API (Producer)
make run-api

# Terminal Two - Private Consumer
make run-consumer
```

**Testing API**

Create a user

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User"}'
```
