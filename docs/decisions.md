# Decisões Técnicas

## Stack escolhida

- Backend em Go
- Framework Gin
- Banco PostgreSQL
- Firebase Auth apenas para autenticação

## Por que essa escolha

- Go é simples, rápido e fácil de manter
- Gin deixa o HTTP leve e direto
- PostgreSQL dá consistência forte e transações confiáveis
- Firebase evita que vocês tenham que construir auth do zero

## Convenções de API

- Prefixo base: `/api/v1`
- JSON em `camelCase`
- Identificadores em path: `trailId`, `missionId`
- Rotas do usuário em `profile` e não em `users/me`

## Decisões de segurança

- O front envia `Authorization: Bearer <idToken>`
- O backend valida o token com Firebase Admin SDK
- O `uid` do Firebase é a identidade real no backend
- Nenhum endpoint de leitura pode expor resposta correta

## Decisões de consistência

- Toda mudança de saldo gera linha em `shell_ledger`
- `users.shell_balance` existe para leitura rápida
- `shell_ledger` é a fonte de verdade financeira
- Submissão deve ser idempotente

## Decisões de produto

- Não migrar o localStorage automaticamente no MVP
- Não criar painel admin no início
- Não usar `users/me` como naming final
- O front passa a consumir API e deixar o estado persistido no backend

## Nomenclatura final de rota

- `GET /api/v1/profile`
- `PATCH /api/v1/profile`
- `PATCH /api/v1/profile/avatar`
- `GET /api/v1/trails`
- `GET /api/v1/trails/:trailId`
- `GET /api/v1/trails/:trailId/missions`
- `GET /api/v1/missions/:missionId`
- `POST /api/v1/missions/:missionId/submissions`
- `GET /api/v1/submissions`
- `GET /api/v1/progress`