## Go‑клиент MAX API

Типизированный Go‑клиент для [бот API платформы MAX](https://dev.max.ru/docs-api). Реализация покрывает эндпоинты, перечисленные в https://dev.max.ru/docs-api:

* Метаданные бота (`GET /me`)
* Управление чатами (`/chats`, `/chats/{chatId}`, `/chats/{chatId}/actions`, `/chats/{chatId}/pin`)
* Участники чатов и роли администраторов (`/members`, `/members/me`, `/members/admins`)
* Подписки и long polling (`/subscriptions`, `/updates`)
* Сессии загрузки файлов (`POST /uploads`)
* Сообщения (`/messages`, `/messages/{messageId}`)
* Метаданные медиа (`/videos/{videoToken}`)
* Ответы на callback (`POST /answers`)

### Установка

```bash
go get github.com/koovod/max-api-go-client
```

### Использование

```go
import (
    "context"
    maxapi "github.com/koovod/max-api-go-client"
)

func main() {
    client, _ := maxapi.NewClient("YOUR_TOKEN")

    me, _ := client.GetMe(context.Background())

    chatID := int64(1234)
    msg, _ := client.SendMessage(
        context.Background(),
        maxapi.SendMessageParams{ChatID: &chatID},
        maxapi.NewMessageBody{Text: "Привет!"}, // отправляем текст и при желании вложения
    )

    _ = msg // обрабатываем ответ
}
```

Основные вспомогательные методы:

| Метод | Описание |
| --- | --- |
| `ListChats`, `GetChat`, `UpdateChat`, `DeleteChat` | Управление метаданными чата и его жизненным циклом. |
| `SendChatAction`, `PinChatMessage`, `UnpinChatMessage` | Контроль пользовательских подсказок в чате. |
| `GetMyChatMember`, `ListChatMembers`, `AddChatMembers`, `RemoveChatMember` | Просмотр и изменение состава участников. |
| `ListSubscriptions`, `CreateSubscription`, `DeleteSubscription`, `GetUpdates` | Управление вебхуками и получение обновлений. |
| `CreateUpload`, `SendMessage`, `EditMessage`, `DeleteMessage`, `GetMessage`, `ListMessages` | Полный набор для работы с сообщениями. |
| `GetVideoMetadata`, `AnswerCallback` | Работа с медиа и ответами на callback‑события. |

Все методы принимают `context.Context`, поэтому можно задавать дедлайны и отмену.

### Тестирование

```bash
cd max-api-go-client
go test ./...
```

Тесты используют `httptest` и не отправляют запросы на реальные сервера MAX.
