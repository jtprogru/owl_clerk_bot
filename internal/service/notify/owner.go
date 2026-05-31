package notify

import (
	"context"
	"fmt"
	"strings"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
	"github.com/jtprogru/owl_clerk_bot/internal/service/sm"
)

// OwnerSender is implemented by the Telegram transport. It knows how to send
// a message with inline action buttons (reply/block/spam/change-category) to
// the owner, attaching the target user's UID to the callback payload.
type OwnerSender interface {
	SendOwnerNotification(ctx context.Context, ownerID, targetUID int64, summary string) error
}

type Owner struct {
	ownerID int64
	sender  OwnerSender
}

func NewOwner(ownerID int64, s OwnerSender) *Owner {
	return &Owner{ownerID: ownerID, sender: s}
}

// Notify composes a summary and sends it to the owner with action buttons.
func (o *Owner) Notify(ctx context.Context, p domain.Profile, data domain.StateData, lastMsg string) error {
	summary := ComposeSummary(p, data, lastMsg)
	return o.sender.SendOwnerNotification(ctx, o.ownerID, p.UID, summary)
}

// ComposeSummary renders a human-friendly summary used in owner notifications
// and in the web UI's conversation card.
func ComposeSummary(p domain.Profile, data domain.StateData, lastMsg string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "👤 %s", p.DisplayName())
	if p.Username != "" {
		fmt.Fprintf(&b, "  @%s", p.Username)
	}
	fmt.Fprintf(&b, "\n📂 Категория: %s\n", p.Category.Title())

	switch p.Category {
	case domain.CategoryHR:
		writeKV(&b, "🏢 Компания", data[sm.KeyCompany])
		writeKV(&b, "💼 Роль", data[sm.KeyRole])
	case domain.CategoryCoop:
		writeKV(&b, "🤝 Проект", data[sm.KeyProject])
	case domain.CategoryFriend:
		writeKV(&b, "🪄 Откуда знакомы", data[sm.KeyOrigin])
	case domain.CategoryMisc:
		writeKV(&b, "📝 Тема", data[sm.KeyTopic])
	}

	contact := p.Contact
	if contact == "" {
		contact = data[sm.KeyContact]
	}
	writeKV(&b, "📞 Контакт", contact)

	if strings.TrimSpace(lastMsg) != "" {
		fmt.Fprintf(&b, "💬 Сообщение: %s\n", lastMsg)
	}
	fmt.Fprintf(&b, "\n#uid_%d", p.UID)
	return b.String()
}

func writeKV(b *strings.Builder, key, val string) {
	val = strings.TrimSpace(val)
	if val == "" {
		return
	}
	fmt.Fprintf(b, "%s: %s\n", key, val)
}
