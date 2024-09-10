package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/frairon/linkbot/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	botApi *tgbotapi.BotAPI
	st     *storage.Storage

	acceptNewUser bool

	mSessions sync.Mutex
	sessions  map[int64]*botSession

	startTime time.Time

	shutdown chan struct{}

	rootState StateBuilder
}

func New(botApi *tgbotapi.BotAPI, st *storage.Storage, rootState StateBuilder) (*Bot, error) {
	return &Bot{
		botApi:    botApi,
		st:        st,
		sessions:  make(map[int64]*botSession),
		shutdown:  make(chan struct{}),
		rootState: rootState,
	}, nil
}

func (b *Bot) getOrCreateSession(ctx context.Context, user *tgbotapi.User, chat *tgbotapi.Chat) (*botSession, error) {
	if chat == nil {
		return nil, fmt.Errorf("chat is nil, cannot create session")
	}
	if user == nil {
		return nil, fmt.Errorf("user is nil, cannot create session")
	}
	b.mSessions.Lock()
	defer b.mSessions.Unlock()

	session := b.sessions[chat.ID]
	if session == nil {
		session = NewSession(b.st, user.ID, chat.ID, b, ctx, b.botApi)
		b.sessions[chat.ID] = session

		// create an initial state and activate
		session.getOrPushCurrentState()
		session.CurrentState().Activate(session)

	}

	if session.chat == nil {
		session.chat = chat
	}

	if session.user == nil {
		session.user = user
	}

	return session, nil
}

func (b *Bot) RootState() State {
	return b.rootState()
}

var (
	CommandReload = tgbotapi.BotCommand{
		Command:     "reload",
		Description: "Reloads the current bot",
	}
	CommandCancel = tgbotapi.BotCommand{
		Command:     "cancel",
		Description: "Cancels the current operation",
	}
	CommandHelp = tgbotapi.BotCommand{
		Command:     "help",
		Description: "Show general help",
	}
	CommandMain = tgbotapi.BotCommand{
		Command:     "home",
		Description: "Goes to the home screen",
	}
	CommandUsers = tgbotapi.BotCommand{
		Command:     "users",
		Description: "Goes to the user management",
	}
)

func (b *Bot) Run(ctx context.Context) error {
	b.startTime = time.Now()
	b.shutdown = make(chan struct{})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botApi.GetUpdatesChan(u)

	// stop the updates
	defer b.botApi.StopReceivingUpdates()

	_, err := b.botApi.Request(tgbotapi.NewSetMyCommands(
		CommandMain,
		CommandUsers,
		CommandCancel,
		// CommandHelp,
		CommandReload))
	if err != nil {
		log.Printf("error setting my commands")
	}

	b.loadSessions(ctx)

	// broadcast shutdown message and store everything
	defer func() {
		for _, session := range b.sessions {
			session.Shutdown()
		}
		b.broadcastToActive("Bot is restarting for maintenance. See you in a few minutes. 🧘")
		b.storeSessions(ctx)
	}()

	sessionStoreTicker := time.NewTicker(60 * time.Second)
	defer sessionStoreTicker.Stop()

	for {
		select {
		case upd, ok := <-updates:
			if !ok {
				return nil
			}

			user := upd.SentFrom()
			if user == nil {
				log.Printf("no sending user - dropping update: %v", upd)
				continue
			}

			if !b.st.UserExists(user.ID) {
				if !b.acceptNewUser {
					log.Printf("user not allowed: %v", user.ID)
					continue
				}

				name := findNameForUser(user)
				log.Printf("Adding new user with %d (%s)", user.ID, name)
				if err := b.st.AddUser(user.ID, name); err != nil {
					log.Printf("Error adding user: %#v: %v", user, err)
					continue
				}
			}

			session, err := b.getOrCreateSession(ctx, user, upd.FromChat())
			if err != nil {
				log.Printf("error handling update %#v: %v", upd, err)
				continue
			}

			if !session.Handle(upd) {
				if upd.Message != nil && upd.Message.Command() != "" {
					command := upd.Message.Command()
					switch command {
					case CommandCancel.Command:
						session.PopState()
					case CommandReload.Command:
						session.ReplaceState(session.CurrentState())
					case CommandHelp.Command:
						session.SendMessage("Help message how to use the bot. TODO.")
					case CommandMain.Command:
						session.ResetToState(b.rootState())
					case CommandUsers.Command:
						session.ResetToState(UsersList())
					default:
						log.Printf("unhandled command: %s", command)
					}
				} else {
					log.Printf("unhandled update: %#v", upd)
				}
			}
		case <-ctx.Done():
			return nil
		case <-b.shutdown:
			log.Printf("bot shutdown initiated")
			return nil
		case <-sessionStoreTicker.C:
			b.storeSessions(ctx)
		}
	}
}

func (b *Bot) foreachSessionAsync(do func(session *botSession)) {
	for _, session := range b.sessions {
		session := session
		go func() {
			do(session)
		}()
	}
}

func (b *Bot) shutdownBot() {
	close(b.shutdown)
}

func findNameForUser(user *tgbotapi.User) string {
	name := user.UserName
	if name == "" {
		name = user.FirstName
	}
	if name == "" {
		name = user.LastName
	}
	if name == "" {
		name = "Unknown"
	}
	return name
}

func (b *Bot) AcceptUsers(ctx context.Context) {
	b.acceptNewUser = true
	go func() {
		select {
		case <-time.After(10 * time.Minute):
			b.acceptNewUser = false
		case <-ctx.Done():
		}
	}()
}

func (b *Bot) storeSessions(ctx context.Context) {
	b.mSessions.Lock()
	defer b.mSessions.Unlock()
	for _, session := range b.sessions {
		data, err := json.Marshal(session.settings)
		if err != nil {
			log.Printf("error marshalling session data: %v", err)
		}
		if err := b.st.StoreSession(session.UserId(), session.ChatId(), session.lastUserAction, string(data)); err != nil {
			log.Printf("error storing session: %v", err)
		}
	}
}

func (b *Bot) loadSessions(ctx context.Context) error {
	b.mSessions.Lock()
	defer b.mSessions.Unlock()

	sessions, err := b.st.UserSessions()
	if err != nil {
		return fmt.Errorf("error loading sessions: %v", err)
	}

	for _, session := range sessions {

		if session.ChatID == 0 || session.UserID == 0 {
			log.Printf("ignoring invalid session: %#v", session)
			continue
		}

		bs := NewSession(b.st, session.UserID, session.ChatID, b, ctx, b.botApi)
		b.sessions[session.ChatID] = bs

		// restore session data
		if session.Data.String != "" {
			var settings SessionSettings
			if err := json.Unmarshal([]byte(session.Data.String), &settings); err == nil {
				bs.settings = &settings
			} else {
				log.Printf("Error restoring session data: %v", err)
			}
		} else {
			bs.settings = &SessionSettings{
				NotifyOnAutoTempChange: true,
			}
		}

		if session.LastUserAction.Valid && time.Since(session.LastUserAction.Time) < time.Hour*24*30 {
			bs.getOrPushCurrentState().Activate(bs)
		} else {
			// initialize to root state
			// TODO: this needs to be some kind of 'init' function instead
			bs.getOrPushCurrentState()
		}

	}

	return nil
}

func (b *Bot) broadcastToActive(message string) {
	b.mSessions.Lock()
	defer b.mSessions.Unlock()

	for _, session := range b.sessions {
		if session.lastUserAction.IsZero() {
			continue
		}
		session.SendMessage(message, SendMessageKeepKeyboard())
	}
}

type BroadcastOptions struct {
	message  string
	newState StateBuilder
}
type BroadcastOption func(opts *BroadcastOptions)

func BroadcastNewState(stateBuilder StateBuilder) BroadcastOption {
	return func(opts *BroadcastOptions) {
		opts.newState = stateBuilder
	}
}

func BroadcastMessage(message string) BroadcastOption {
	return func(opts *BroadcastOptions) {
		opts.message = message
	}
}

func BroadcastMessagef(format string, args ...interface{}) BroadcastOption {
	return BroadcastMessage(fmt.Sprintf(format, args...))
}

func (b *Bot) broadcast(opts ...BroadcastOption) {
	b.mSessions.Lock()
	defer b.mSessions.Unlock()

	for _, session := range b.sessions {
		var options BroadcastOptions

		for _, opt := range opts {
			opt(&options)
		}

		session := session
		go func() {
			if options.message != "" {
				session.SendMessage(options.message, SendMessageKeepKeyboard())
			}
			if options.newState != nil {
				session.ResetToState(options.newState())
			}
		}()

	}
}
