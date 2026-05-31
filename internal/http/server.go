package http

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	stdhttp "net/http"
	"strconv"
	"time"

	"github.com/jtprogru/owl_clerk_bot/internal/domain"
	"github.com/jtprogru/owl_clerk_bot/internal/storage/sqlite"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

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

type OutboundService interface {
	Reply(ctx context.Context, uid int64, text string) error
}

type Server struct {
	addr     string
	user     string
	pass     string
	log      *slog.Logger
	profiles ProfileManager
	messages MessageReader
	outbound OutboundService

	dashboardTpl *template.Template
	listTpl      *template.Template
	cardTpl      *template.Template
}

func New(
	addr, user, pass string,
	log *slog.Logger,
	profiles ProfileManager,
	messages MessageReader,
	outbound OutboundService,
) (*Server, error) {
	s := &Server{
		addr:     addr,
		user:     user,
		pass:     pass,
		log:      log,
		profiles: profiles,
		messages: messages,
		outbound: outbound,
	}
	var err error
	if s.dashboardTpl, err = parsePage("dashboard.html"); err != nil {
		return nil, err
	}
	if s.listTpl, err = parsePage("list.html"); err != nil {
		return nil, err
	}
	if s.cardTpl, err = parsePage("card.html"); err != nil {
		return nil, err
	}
	return s, nil
}

func parsePage(name string) (*template.Template, error) {
	return template.ParseFS(templatesFS, "templates/layout.html", "templates/"+name)
}

func (s *Server) Run(ctx context.Context) error {
	mux := stdhttp.NewServeMux()
	mux.HandleFunc("GET /", s.handleDashboard)
	mux.HandleFunc("GET /conversations", s.handleList)
	mux.HandleFunc("GET /conversations/{uid}", s.handleCard)
	mux.HandleFunc("POST /conversations/{uid}/reply", s.handleReply)
	mux.HandleFunc("POST /conversations/{uid}/category", s.handleSetCategory)
	mux.HandleFunc("POST /conversations/{uid}/block", s.handleSetBlocked)

	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return fmt.Errorf("static fs: %w", err)
	}
	mux.Handle("GET /static/", stdhttp.StripPrefix("/static/", stdhttp.FileServer(stdhttp.FS(staticSub))))

	srv := &stdhttp.Server{
		Addr:              s.addr,
		Handler:           basicAuth(s.user, s.pass, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		s.log.Info("web server starting", "addr", s.addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

// --- handlers ---

func (s *Server) handleDashboard(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()
	counts, err := s.profiles.CountByCategory(ctx)
	if err != nil {
		s.serverError(w, err)
		return
	}
	week, err := s.messages.CountSince(ctx, 7)
	if err != nil {
		s.serverError(w, err)
		return
	}
	month, err := s.messages.CountSince(ctx, 30)
	if err != nil {
		s.serverError(w, err)
		return
	}
	s.render(w, s.dashboardTpl, map[string]any{
		"Title":  "дашборд",
		"Counts": counts,
		"Week":   week,
		"Month":  month,
	})
}

func (s *Server) handleList(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()
	q := r.URL.Query()
	f := sqlite.ListFilter{Limit: 100}
	if c := q.Get("category"); c != "" {
		cat := domain.Category(c)
		if cat.Valid() {
			f.Category = cat
		}
	}
	if q.Get("blocked") == "1" {
		b := true
		f.Blocked = &b
	}
	profiles, err := s.profiles.List(ctx, f)
	if err != nil {
		s.serverError(w, err)
		return
	}
	s.render(w, s.listTpl, map[string]any{
		"Title":    "диалоги",
		"Profiles": profiles,
		"Filter":   filterView{Category: f.Category, Blocked: f.Blocked != nil && *f.Blocked},
	})
}

type filterView struct {
	Category domain.Category
	Blocked  bool
}

func (s *Server) handleCard(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()
	uid, err := uidFromPath(r)
	if err != nil {
		stdhttp.Error(w, "bad uid", stdhttp.StatusBadRequest)
		return
	}
	p, err := s.profiles.Get(ctx, uid)
	if errors.Is(err, domain.ErrNotFound) {
		stdhttp.NotFound(w, r)
		return
	}
	if err != nil {
		s.serverError(w, err)
		return
	}
	msgs, err := s.messages.ListByUID(ctx, uid, 200)
	if err != nil {
		s.serverError(w, err)
		return
	}
	s.render(w, s.cardTpl, map[string]any{
		"Title":    p.DisplayName(),
		"Profile":  p,
		"Messages": msgs,
	})
}

func (s *Server) handleReply(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	uid, err := uidFromPath(r)
	if err != nil {
		stdhttp.Error(w, "bad uid", stdhttp.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		stdhttp.Error(w, "bad form", stdhttp.StatusBadRequest)
		return
	}
	text := r.PostFormValue("text")
	if text == "" {
		stdhttp.Error(w, "empty text", stdhttp.StatusBadRequest)
		return
	}
	if err := s.outbound.Reply(r.Context(), uid, text); err != nil {
		s.serverError(w, err)
		return
	}
	stdhttp.Redirect(w, r, fmt.Sprintf("/conversations/%d", uid), stdhttp.StatusSeeOther)
}

func (s *Server) handleSetCategory(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	uid, err := uidFromPath(r)
	if err != nil {
		stdhttp.Error(w, "bad uid", stdhttp.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		stdhttp.Error(w, "bad form", stdhttp.StatusBadRequest)
		return
	}
	cat := domain.Category(r.PostFormValue("category"))
	if !cat.Valid() {
		stdhttp.Error(w, "bad category", stdhttp.StatusBadRequest)
		return
	}
	if err := s.profiles.SetCategory(r.Context(), uid, cat); err != nil {
		s.serverError(w, err)
		return
	}
	stdhttp.Redirect(w, r, fmt.Sprintf("/conversations/%d", uid), stdhttp.StatusSeeOther)
}

func (s *Server) handleSetBlocked(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	uid, err := uidFromPath(r)
	if err != nil {
		stdhttp.Error(w, "bad uid", stdhttp.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		stdhttp.Error(w, "bad form", stdhttp.StatusBadRequest)
		return
	}
	blocked := r.PostFormValue("blocked") == "1"
	if err := s.profiles.SetBlocked(r.Context(), uid, blocked); err != nil {
		s.serverError(w, err)
		return
	}
	stdhttp.Redirect(w, r, fmt.Sprintf("/conversations/%d", uid), stdhttp.StatusSeeOther)
}

// --- helpers ---

func uidFromPath(r *stdhttp.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("uid"), 10, 64)
}

func (s *Server) render(w stdhttp.ResponseWriter, tpl *template.Template, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		s.log.Error("render", "err", err)
	}
}

func (s *Server) serverError(w stdhttp.ResponseWriter, err error) {
	s.log.Error("server error", "err", err)
	stdhttp.Error(w, "internal error", stdhttp.StatusInternalServerError)
}
