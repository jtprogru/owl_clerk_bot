package tg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
	"github.com/jtprogru/owl_clerk_bot/internal/service/intake"
	"github.com/jtprogru/owl_clerk_bot/internal/storage/sqlite"
	tele "gopkg.in/telebot.v3"
)

// ProfileManager is the slice of the profile store the transport needs.
type ProfileManager interface {
	Get(ctx context.Context, uid int64) (domain.Profile, error)
	List(ctx context.Context, f sqlite.ListFilter) ([]domain.Profile, error)
	SetCategory(ctx context.Context, uid int64, c domain.Category) error
	SetBlocked(ctx context.Context, uid int64, blocked bool) error
	CountByCategory(ctx context.Context) ([]sqlite.CategoryCount, error)
}

type MessageReader interface {
	ListByUID(ctx context.Context, uid int64, limit int) ([]domain.Message, error)
	CountSince(ctx context.Context, days int) (int, error)
}

type IntakeService interface {
	Handle(ctx context.Context, in intake.Incoming) (intake.Response, error)
}

type OutboundService interface {
	Reply(ctx context.Context, uid int64, text string) error
}

// Handlers wires Telegram callbacks to the service layer.
type Handlers struct {
	ownerID  int64
	log      *slog.Logger
	intake   IntakeService
	outbound OutboundService
	profiles ProfileManager
	messages MessageReader
}

func NewHandlers(
	ownerID int64,
	log *slog.Logger,
	in IntakeService,
	out OutboundService,
	profiles ProfileManager,
	messages MessageReader,
) *Handlers {
	return &Handlers{
		ownerID:  ownerID,
		log:      log,
		intake:   in,
		outbound: out,
		profiles: profiles,
		messages: messages,
	}
}

func (h *Handlers) Register(b *tele.Bot) {
	b.Use(recoverMW(h.log))

	owner := b.Group()
	owner.Use(ownerOnly(h.ownerID))
	owner.Handle("/start", h.cmdStart)
	owner.Handle("/list", h.cmdList)
	owner.Handle("/show", h.cmdShow)
	owner.Handle("/reply", h.cmdReply)
	owner.Handle("/block", h.cmdBlock)
	owner.Handle("/unblock", h.cmdUnblock)
	owner.Handle("/category", h.cmdCategory)
	owner.Handle("/stats", h.cmdStats)

	// Inline-кнопки уведомлений (только владелец).
	owner.Handle(&tele.Btn{Unique: btnBlock}, h.onBtnBlock)
	owner.Handle(&tele.Btn{Unique: btnSpam}, h.onBtnSpam)
	owner.Handle(&tele.Btn{Unique: btnCatPick}, h.onBtnCatPick)
	owner.Handle(&tele.Btn{Unique: btnCatSet}, h.onBtnCatSet)

	// Текст — общий, маршрутизация по sender внутри.
	b.Handle(tele.OnText, h.onText)
}

// onText is the main router for plain text messages.
func (h *Handlers) onText(c tele.Context) error {
	if c.Sender() != nil && c.Sender().ID == h.ownerID {
		return h.onOwnerText(c)
	}
	return h.onUserText(c)
}

func (h *Handlers) onUserText(c tele.Context) error {
	ctx := context.Background()
	in := intake.Incoming{
		Profile: profileFromCtx(c),
		Text:    c.Text(),
		TgMsgID: int64(c.Message().ID),
	}
	resp, err := h.intake.Handle(ctx, in)
	if err != nil {
		h.log.Error("intake handle", "uid", in.Profile.UID, "err", err)
		return c.Send("Что-то пошло не так. Попробуй ещё раз чуть позже.")
	}
	if resp.Blocked || resp.Reply.Text == "" {
		return nil
	}
	if kbd := userReplyKeyboard(resp.Reply.Keyboard); kbd != nil {
		return c.Send(resp.Reply.Text, kbd)
	}
	return c.Send(resp.Reply.Text, clearReplyKeyboard())
}

// onOwnerText handles plain text from the owner. The owner-side flow is:
//   - reply-to a notification → route as outbound to target uid
//   - anything else → hint
func (h *Handlers) onOwnerText(c tele.Context) error {
	if r := c.Message().ReplyTo; r != nil {
		if uid, ok := extractTargetUID(r.Text); ok {
			if err := h.outbound.Reply(context.Background(), uid, c.Text()); err != nil {
				h.log.Error("owner reply via reply-to", "uid", uid, "err", err)
				return c.Reply("Не удалось отправить.")
			}
			return c.Reply(fmt.Sprintf("Отправил → uid %d", uid))
		}
	}
	return c.Reply("Не понял. Используй /reply <uid> <текст> или ответь на уведомление.")
}

// --- Commands ---

func (h *Handlers) cmdStart(c tele.Context) error {
	return c.Send(strings.Join([]string{
		"Я бот-секретарь. Доступные команды:",
		"/list [hr|coop|friend|misc|spam] — последние диалоги",
		"/show <uid> — карточка диалога",
		"/reply <uid> <текст> — ответить",
		"/block <uid> | /unblock <uid>",
		"/category <uid> <hr|coop|friend|misc|spam>",
		"/stats — счётчики за неделю/месяц",
	}, "\n"))
}

func (h *Handlers) cmdList(c tele.Context) error {
	ctx := context.Background()
	f := sqlite.ListFilter{Limit: 20}
	if arg := strings.TrimSpace(c.Message().Payload); arg != "" {
		cat := domain.Category(arg)
		if !cat.Valid() {
			return c.Reply("Неизвестная категория. Допустимы: hr, coop, friend, misc, spam.")
		}
		f.Category = cat
	}
	profiles, err := h.profiles.List(ctx, f)
	if err != nil {
		h.log.Error("list profiles", "err", err)
		return c.Reply("Ошибка чтения.")
	}
	if len(profiles) == 0 {
		return c.Reply("Пусто.")
	}
	var b strings.Builder
	for _, p := range profiles {
		flag := ""
		if p.IsBlocked {
			flag = " 🚫"
		}
		fmt.Fprintf(&b, "• [%d] %s — %s%s\n", p.UID, p.DisplayName(), p.Category.Title(), flag)
	}
	return c.Reply(b.String())
}

func (h *Handlers) cmdShow(c tele.Context) error {
	ctx := context.Background()
	uid, err := parseUID(c.Message().Payload)
	if err != nil {
		return c.Reply("Использование: /show <uid>")
	}
	p, err := h.profiles.Get(ctx, uid)
	if errors.Is(err, domain.ErrNotFound) {
		return c.Reply("Профиль не найден.")
	}
	if err != nil {
		return c.Reply("Ошибка чтения.")
	}
	msgs, err := h.messages.ListByUID(ctx, uid, 20)
	if err != nil {
		return c.Reply("Ошибка чтения сообщений.")
	}

	var b strings.Builder
	fmt.Fprintf(&b, "👤 %s (@%s)\n📂 %s%s\n📞 %s\n\n",
		p.DisplayName(), p.Username, p.Category.Title(),
		blockedFlag(p.IsBlocked), valueOr(p.Contact, "—"))
	for _, m := range msgs {
		arrow := "→"
		if m.Direction == domain.DirectionIn {
			arrow = "←"
		}
		fmt.Fprintf(&b, "%s %s: %s\n", arrow, m.At.Format("02.01 15:04"), m.Text)
	}
	fmt.Fprintf(&b, "\n#uid_%d", uid)
	return c.Reply(b.String())
}

func (h *Handlers) cmdReply(c tele.Context) error {
	payload := strings.TrimSpace(c.Message().Payload)
	parts := strings.SplitN(payload, " ", 2)
	if len(parts) < 2 {
		return c.Reply("Использование: /reply <uid> <текст>")
	}
	uid, err := parseUID(parts[0])
	if err != nil {
		return c.Reply("Неверный uid.")
	}
	if err := h.outbound.Reply(context.Background(), uid, parts[1]); err != nil {
		h.log.Error("cmd reply", "uid", uid, "err", err)
		return c.Reply("Не удалось отправить.")
	}
	return c.Reply(fmt.Sprintf("Отправил → uid %d", uid))
}

func (h *Handlers) cmdBlock(c tele.Context) error   { return h.setBlocked(c, true) }
func (h *Handlers) cmdUnblock(c tele.Context) error { return h.setBlocked(c, false) }

func (h *Handlers) setBlocked(c tele.Context, blocked bool) error {
	uid, err := parseUID(c.Message().Payload)
	if err != nil {
		return c.Reply("Использование: /block <uid> или /unblock <uid>")
	}
	if err := h.profiles.SetBlocked(context.Background(), uid, blocked); err != nil {
		return c.Reply("Ошибка.")
	}
	state := "разблокирован"
	if blocked {
		state = "заблокирован"
	}
	return c.Reply(fmt.Sprintf("uid %d %s.", uid, state))
}

func (h *Handlers) cmdCategory(c tele.Context) error {
	parts := strings.Fields(c.Message().Payload)
	if len(parts) != 2 {
		return c.Reply("Использование: /category <uid> <hr|coop|friend|misc|spam>")
	}
	uid, err := parseUID(parts[0])
	if err != nil {
		return c.Reply("Неверный uid.")
	}
	cat := domain.Category(parts[1])
	if !cat.Valid() {
		return c.Reply("Неизвестная категория.")
	}
	if err := h.profiles.SetCategory(context.Background(), uid, cat); err != nil {
		return c.Reply("Ошибка.")
	}
	return c.Reply(fmt.Sprintf("uid %d → %s", uid, cat.Title()))
}

func (h *Handlers) cmdStats(c tele.Context) error {
	ctx := context.Background()
	counts, err := h.profiles.CountByCategory(ctx)
	if err != nil {
		return c.Reply("Ошибка.")
	}
	week, _ := h.messages.CountSince(ctx, 7)
	month, _ := h.messages.CountSince(ctx, 30)

	var b strings.Builder
	b.WriteString("📊 Профили по категориям:\n")
	for _, c := range counts {
		fmt.Fprintf(&b, "• %s: %d\n", c.Category.Title(), c.Count)
	}
	fmt.Fprintf(&b, "\n💬 Сообщений за 7 дней: %d\n💬 За 30 дней: %d", week, month)
	return c.Reply(b.String())
}

// --- Callback buttons ---

func (h *Handlers) onBtnBlock(c tele.Context) error {
	uid, err := parseUID(c.Callback().Data)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Bad uid"})
	}
	if err := h.profiles.SetBlocked(context.Background(), uid, true); err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка"})
	}
	return c.Respond(&tele.CallbackResponse{Text: fmt.Sprintf("uid %d заблокирован", uid)})
}

func (h *Handlers) onBtnSpam(c tele.Context) error {
	ctx := context.Background()
	uid, err := parseUID(c.Callback().Data)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Bad uid"})
	}
	if err := h.profiles.SetCategory(ctx, uid, domain.CategorySpam); err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка"})
	}
	_ = h.profiles.SetBlocked(ctx, uid, true)
	return c.Respond(&tele.CallbackResponse{Text: fmt.Sprintf("uid %d → спам", uid)})
}

func (h *Handlers) onBtnCatPick(c tele.Context) error {
	uid, err := parseUID(c.Callback().Data)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Bad uid"})
	}
	if err := c.Edit(c.Message().Text, categoryPickerKeyboard(uid)); err != nil {
		h.log.Error("edit markup", "err", err)
	}
	return c.Respond()
}

func (h *Handlers) onBtnCatSet(c tele.Context) error {
	parts := strings.SplitN(c.Callback().Data, ":", 2)
	if len(parts) != 2 {
		return c.Respond(&tele.CallbackResponse{Text: "Bad data"})
	}
	uid, err := parseUID(parts[0])
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Bad uid"})
	}
	cat := domain.Category(parts[1])
	if !cat.Valid() {
		return c.Respond(&tele.CallbackResponse{Text: "Bad cat"})
	}
	if err := h.profiles.SetCategory(context.Background(), uid, cat); err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка"})
	}
	if err := c.Edit(c.Message().Text, ownerActionsKeyboard(uid)); err != nil {
		h.log.Error("edit markup back", "err", err)
	}
	return c.Respond(&tele.CallbackResponse{Text: fmt.Sprintf("→ %s", cat.Title())})
}

// --- helpers ---

func parseUID(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty uid")
	}
	return strconv.ParseInt(s, 10, 64)
}

func blockedFlag(b bool) string {
	if b {
		return " 🚫"
	}
	return ""
}

func valueOr(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}
