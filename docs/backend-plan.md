# Plano do Backend

## Objetivo

Substituir o estado persistido no `localStorage` por um backend seguro, com autenticação via Firebase, persistência em PostgreSQL e regras de negócio centralizadas.

## Escopo

- Autenticação via Firebase Auth
- Perfil do usuário
- Trilhas e missões
- Submissões com correção no backend
- Saldo de conchas com auditoria
- Progresso derivado de histórico

## Regras críticas

### 1. Ledger de conchas

O saldo do usuário não deve ser atualizado de forma solta. Toda alteração precisa gerar um registro em `shell_ledger`.

- `users.shell_balance` é projeção para leitura rápida
- `shell_ledger` é a fonte de verdade financeira
- Toda recompensa ou ajuste precisa ser transacional

### 2. Idempotência em submissão

O endpoint de submissão precisa aceitar uma chave de idempotência.

- Requisições repetidas não podem gerar recompensa duplicada
- A mesma chave para o mesmo usuário deve retornar o resultado anterior
- Isso evita duplicação por retry do front ou duplo clique

### 3. Respostas corretas nunca expostas

Nenhum endpoint de leitura pode retornar gabarito, hash de resposta ou qualquer equivalente.

- A correção acontece apenas no backend
- O front recebe apenas conteúdo seguro para exibição
- Missão pode expor enunciado, opções e metadados, mas nunca a resposta correta

### 4. Progresso derivado

Os campos que hoje vivem no `localStorage` como estado principal passam a ser derivados.

- `completed`
- `niveis_concluidos`
- `conchas_por_missao`

Esses dados devem ser calculados a partir de submissões e ledger, não gravados como verdade principal no perfil.

## API final

### Perfil

- `GET /api/v1/profile`
- `PATCH /api/v1/profile`
- `PATCH /api/v1/profile/avatar`

### Conteúdo

- `GET /api/v1/trails`
- `GET /api/v1/trails/:trailId`
- `GET /api/v1/trails/:trailId/missions`
- `GET /api/v1/missions/:missionId`

### Progresso e submissões

- `POST /api/v1/missions/:missionId/submissions`
- `GET /api/v1/submissions`
- `GET /api/v1/progress`

## Fluxo de auth

1. O front autentica com Firebase SDK.
2. O front envia `Authorization: Bearer <idToken>` nas requests.
3. O backend valida o token com Firebase Admin SDK.
4. O backend extrai o `uid` e usa isso como identidade do usuário.

## Fases de implementação

### Fase 1

- Setup do projeto Go
- Gin e middleware de auth
- PostgreSQL e migrations
- Healthcheck e estrutura base

### Fase 2

- Modelos de dados
- Leitura de trilhas e missões
- Perfil do usuário

### Fase 3

- Submissão com correção
- Idempotência
- Ledger de conchas

### Fase 4

- Progresso derivado
- Cosméticos e desbloqueios
- Testes de concorrência

## Critério de sucesso

- Não existe duplicação de recompensa por retry
- Nenhuma resposta correta vaza em leitura
- Progresso é reconstruível a partir do histórico
- Ledger e saldo batem sempre