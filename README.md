# abstract-api

Backend em Go + Gin para o produto Abstract.io — API central para autenticação, perfil, progresso, submissões e saldo de conchas.

## Como começar

1. Crie um `.env` a partir de `.env.example` e preencha `DATABASE_URL`.
2. Ajuste `ADDR` para mudar a porta do servidor (padrão `:8080`).
3. Execute a API com `go run ./cmd/api`.
4. Verifique o healthcheck em `GET /healthz` ou `GET /api/v1/health`.

## Estrutura

- `cmd/api`: ponto de entrada
- `internal/server`: router e rotas base
