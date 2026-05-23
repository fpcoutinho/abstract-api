# Modelo de Domínio

## Entidades

### users

- `uid` primary key vindo do Firebase
- `name`
- `gender`
- `shell_balance`
- `avatar_idx`
- `active_frame`
- `active_accessory`
- `created_at`
- `updated_at`

### trails

- `id`
- `slug`
- `title`
- `description`
- `order_index`
- `published`

### missions

- `id`
- `trail_id`
- `slug`
- `title`
- `content_md`
- `reward_shells`
- `order_index`
- `published`
- `answer_key_hash`

### user_submissions

- `id`
- `uid`
- `mission_id`
- `answers_json`
- `is_correct`
- `earned_shells`
- `idempotency_key`
- `completed_at`
- `created_at`

### shell_ledger

- `id`
- `uid`
- `submission_id`
- `delta`
- `reason`
- `balance_before`
- `balance_after`
- `created_at`

### user_unlocks

- `id`
- `uid`
- `item_type`
- `item_code`
- `unlocked_at`

## Relações

- `trails` possui muitas `missions`
- `users` possui muitas `user_submissions`
- `users` possui muitos `shell_ledger`
- `user_submissions` referencia uma `mission`
- `shell_ledger` referencia uma `submission`
- `user_unlocks` pertence a um `user`

## Regras de consistência

- `shell_balance` é projeção rápida do saldo
- `shell_ledger` é a fonte de verdade do histórico financeiro
- `user_submissions` precisa de chave de idempotência única por usuário
- Uma missão correta só recompensa uma vez, mesmo que seja reenviada
- Campos de progresso devem ser derivados do histórico, não copiados como estado principal

## Fluxo de submissão

1. O front envia a submissão com `idempotency_key`.
2. O backend valida o token Firebase e identifica o usuário.
3. O backend busca a missão.
4. O backend corrige a resposta.
5. O backend grava submissão, ledger e saldo em uma única transação.
6. O backend responde com resultado seguro para o front.

## Dados que não devem ser expostos

- `answer_key_hash`
- gabarito da missão
- qualquer estrutura que permita reconstruir a resposta correta no front