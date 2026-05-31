# Technical specification — owl_clerk_bot v1

## Назначение

Личный бот-секретарь для Telegram, прикреплённый к аккаунту владельца
(@jtprogru). Принимает обращения от незнакомцев, проводит короткий
структурированный опрос, сохраняет результат в SQLite и уведомляет владельца
в личке. Управление — через команды в чате с ботом, inline-кнопки в
уведомлениях и веб-морду.

## Функциональные требования

### Сценарий пользователя (незнакомец)
1. Greeting — приветствие от бота.
2. AskCategory — выбор категории кнопками: HR, Сотрудничество, Дружба,
   Разное.
3. Категорийная ветка — 1 уточняющий вопрос:
   - HR: компания → роль
   - Сотрудничество: проект
   - Дружба: где познакомились
   - Разное: тема
4. AskContact — оставить контакт.
5. Done — подтверждение, что сводка передана владельцу.
6. Любое сообщение после Done — молча сохраняется и нотифицирует владельца.

### Владелец
- Команды в ЛС бота: `/start`, `/list [cat]`, `/show <uid>`, `/reply
  <uid> <текст>`, `/block <uid>`, `/unblock <uid>`,
  `/category <uid> <cat>`, `/stats`.
- Reply-to уведомления (`#uid_N` в тексте) — выступает как `/reply`.
- Inline-кнопки в уведомлении: 🚫 Заблок, 🗑 Спам, 🔁 Сменить категорию.
- Веб-морда: дашборд, список диалогов с фильтрами, карточка диалога с
  историей и формой ответа.

## Архитектура

```
cmd/app              entrypoint, сборка зависимостей, graceful shutdown
internal/config      .env → config
internal/domain      чистые типы (Profile, Message, Category, StateID, …)
internal/storage/sqlite  миграции (embed) + ProfileRepo/MessageRepo/StateRepo
internal/service/sm       FSM, состояния, registry
internal/service/intake   координатор приёма входящих
internal/service/notify   формирование сводки + отправка владельцу
internal/service/outbound out-of-band ответы (команды, веб)
internal/transport/tg     Bot, Client, Handlers, Keyboards
internal/http             Web UI (ServeMux + html/template + embed)
```

## Хранилище (SQLite)

См. `internal/storage/sqlite/migrations/0001_init.sql`. Таблицы:
- `profiles` — UID, имя, username, категория, контакт, флаг блока, метки
  времени.
- `messages` — лог in/out сообщений с привязкой к UID.
- `conversation_states` — текущее состояние FSM + собранные ответы (JSON).

## Конфигурация (`.env`)

| Переменная       | Назначение                                       |
|------------------|--------------------------------------------------|
| OWL_BOT_TOKEN    | токен бота от @BotFather                         |
| OWL_OWNER_ID     | telegram user id владельца (форвард + ACL)       |
| OWL_DEBUG        | подробные логи + telebot verbose                 |
| OWL_SERVE_HOST   | хост веб-морды                                   |
| OWL_SERVE_PORT   | порт веб-морды                                   |
| OWL_DB_PATH      | путь к SQLite-файлу                              |
| OWL_WEB_USER     | пользователь HTTP basic auth                     |
| OWL_WEB_PASS     | пароль HTTP basic auth                           |

## Стек

- Go 1.25 (slog, ServeMux pattern matching, embed)
- `gopkg.in/telebot.v3` — Telegram long polling
- `modernc.org/sqlite` — pure-Go SQLite (без CGO)
- stdlib `net/http` + `html/template`

## Verification

- `task build:bin` — собирается без ошибок;
- `task vet` — пусто;
- `task test` — все тесты зелёные (sm, intake, storage/sqlite);
- ручной прогон пользовательского сценария по всем 4 веткам;
- проверка owner-команд и inline-кнопок;
- веб-морда: дашборд, фильтр по категории, ответ из карточки.

## Бэклог (v2+)

- Anti-spam rate-limit (1 опрос в сутки на UID);
- «Режим отпуска» (автоответ «вернусь Х»);
- Экспорт диалога в Markdown/JSON;
- Telegram OAuth для веб-морды вместо basic auth;
- Метрики (Prometheus) для дашборда.
