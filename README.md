# Kita Go Service

Microservice for GO

## Getting Started

If you use nix: `nix-shell` or `nix-shell shell.nix`

**Start infrastructure**

You can skip straight to `Testing API` step if you dont need to run the apps independently

```bash
make docker-up
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

## Schema Evolution

Using [protobuf](https://protobuf.dev/overview/) to manage event schemas. Proto files are located in `proto`, use `make proto` to generate code which will will output to `internal/events/proto`
