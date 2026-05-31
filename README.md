# owl_clerk_bot

Личный бот-секретарь для Telegram. Принимает обращения от незнакомых людей,
проводит короткий опрос (HR / Сотрудничество / Дружба / Разное), сохраняет
профиль и контакты в SQLite и форвардит сводку владельцу. Управление —
через команды в личке, inline-кнопки в уведомлениях или веб-морду.

## Запуск

Создать `.env` (в репо не коммитится):

```env
OWL_BOT_TOKEN=1234567890:AA...           # токен бота от @BotFather
OWL_OWNER_ID=123456789                   # твой telegram user id
OWL_DEBUG=false
OWL_SERVE_HOST=127.0.0.1
OWL_SERVE_PORT=8000
OWL_DB_PATH=./data/owl.db
OWL_WEB_USER=admin
OWL_WEB_PASS=change-me
```

Локально:

```shell
task run:bin
```

Доступные команды:

```shell
task --list
```

## Стек

- Go 1.25+
- `gopkg.in/telebot.v3` — Telegram long polling
- `modernc.org/sqlite` — pure-Go SQLite (без CGO)
- `log/slog` (stdlib)
- `net/http` + `html/template` для веб-морды

## Структура

```
cmd/app                — entrypoint
internal/config        — env-конфиг
internal/domain        — типы (Profile, Message, Category, State)
internal/storage/sqlite — миграции и репозитории
internal/service/sm    — конечный автомат опроса
internal/service/intake   — координатор приёма
internal/service/notify   — уведомления владельцу
internal/service/outbound — исходящие ответы
internal/transport/tg  — Telegram handler'ы
internal/http          — web UI
```

См. также [TS.md](./TS.md) и план в `~/.claude/plans/`.
