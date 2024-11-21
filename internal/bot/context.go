package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/frairon/linkbot/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GlobalMessageHandler func(bs *BotSession, message *tgbotapi.Message) bool

type SessionSettings struct {
	PauseAllNotifications bool

	NotifyOnBell             bool
	NotifyOnButtonTempChange bool
	NotifyOnAutoTempChange   bool
	NotifyOnMovement         bool
	NotifyOnEnterLeave       bool
	NotifyOnTempChange       bool

	SingleMovementAlert   bool
	SingleEnterLeaveAlert bool
}

type BotSession struct {
	st *storage.Storage

	userId int64
	chatId int64

	user *tgbotapi.User
	chat *tgbotapi.Chat

	lastUserAction time.Time

	states []State

	globalMessageHandler []GlobalMessageHandler

	bot    *Bot
	botCtx context.Context

	botApi *tgbotapi.BotAPI

	sessionCommandHandlers map[string]CommandHandler

	settings *SessionSettings
}

func NewSession(st *storage.Storage, userId int64, chatId int64, bot *Bot, botCtx context.Context, botApi *tgbotapi.BotAPI) *BotSession {
	return &BotSession{
		st:                     st,
		userId:                 userId,
		chatId:                 chatId,
		bot:                    bot,
		botCtx:                 botCtx,
		botApi:                 botApi,
		sessionCommandHandlers: make(map[string]CommandHandler),
		globalMessageHandler: []GlobalMessageHandler{
			GlobalHomeHandler,
		},
		settings: &SessionSettings{},
	}
}

func (bs *BotSession) Settings() *SessionSettings {
	return bs.settings
}

func (bs *BotSession) getOrPushCurrentState() State {
	if len(bs.states) == 0 {
		bs.states = []State{bs.bot.RootState()}
	}

	return bs.states[len(bs.states)-1]
}

func (bs *BotSession) Handle(update tgbotapi.Update) bool {
	curState := bs.getOrPushCurrentState()

	bs.lastUserAction = time.Now()

	switch {
	case update.Message != nil:

		// if the message is a command, try to handle that instead.
		// First the current stae, then the context
		if cmd := update.Message.CommandWithAt(); cmd != "" {
			args := strings.Split(update.Message.CommandArguments(), " ")
			if curState.HandleCommand(bs, cmd, args...) {
				return true
			}
			return bs.handleCommand(cmd, args)
		}

		for _, handler := range bs.globalMessageHandler {
			if handler(bs, update.Message) {
				return true
			}
		}

		return curState.HandleMessage(bs, update.Message)
	case update.CallbackQuery != nil:

		if curState.HandleCallbackQuery(bs, update.CallbackQuery) {
			return true
		} else {
			return bs.removeExpiredCallback(update.CallbackQuery)
		}

	default:
		log.Printf("unhandled update: %#v", update)
	}
	return false
}

func (bs *BotSession) removeExpiredCallback(query *tgbotapi.CallbackQuery) bool {
	alert := tgbotapi.NewCallbackWithAlert(query.InlineMessageID, "message expired, buttons disabled")
	alert.CallbackQueryID = query.ID

	if query.Message != nil {
		bs.RemoveKeyboardForMessage(query.Message.MessageID)
	}
	_, err := bs.botApi.Request(alert)
	if err != nil {
		bs.SendError(err)
	}
	return true
}

func (bs *BotSession) RemoveKeyboardForMessage(messageId int) {
	// construct an update reply-markup message manually, because we need to set
	// the ReplyMarkup to nil, which is not supported by the library
	bs.botApi.Request(tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      bs.chatId,
			MessageID:   messageId,
			ReplyMarkup: nil,
		},
	})
}

func (bs *BotSession) handleCommand(command string, args []string) bool {
	switch command {
	case CommandCancel.Command:
		bs.PopState()
		return true
	}

	for _, handler := range bs.sessionCommandHandlers {
		if handler.Handle(bs, command, args...) {
			return true
		}
	}

	return false
}

func (bs *BotSession) SetCommandHandler(name string, handler CommandHandler) {
	bs.sessionCommandHandlers[name] = handler
}

func (bs *BotSession) PushState(state State) {
	if len(bs.states) > 0 {
		bs.CurrentState().BeforeLeave(bs)
	}
	bs.states = append(bs.states, state)
	state.Activate(bs)
}

func (bs *BotSession) PopState() {
	if len(bs.states) == 0 {
		return
	}

	bs.CurrentState().BeforeLeave(bs)

	bs.states = bs.states[:len(bs.states)-1]

	curState := bs.getOrPushCurrentState()

	curState.Return(bs)
}

func (bs *BotSession) DropStates(n int) {
	if len(bs.states) > n {
		bs.states = bs.states[:len(bs.states)-n]
	} else {
		bs.states = nil
	}
	bs.getOrPushCurrentState().Return(bs)
}

func (bs *BotSession) CurrentState() State {
	if len(bs.states) == 0 {
		return nil
	}
	return bs.states[len(bs.states)-1]
}

func (bs *BotSession) ReplaceState(state State) {
	if len(bs.states) == 0 {
		return
	}

	bs.states[len(bs.states)-1] = state
	state.Activate(bs)
}

func (bs *BotSession) ResetToState(state State) {
	bs.states = nil
	bs.PushState(state)
}

func (bs *BotSession) Storage() *storage.Storage {
	return bs.st
}

func (bs *BotSession) UserId() int64 {
	return bs.userId
}

func (bs *BotSession) ChatId() int64 {
	return bs.chatId
}

func (bc *BotSession) SendMessageWithCommands(text string, replyCommands ButtonKeyboard, opts ...SendMessageOption) int {
	msg := tgbotapi.NewMessage(bc.ChatId(), text)
	msg.ParseMode = "html"

	options := &sendMessageOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if replyCommands != nil {
		keyboard := tgbotapi.ReplyKeyboardMarkup{
			ResizeKeyboard: true,
		}
		for _, row := range replyCommands {
			// rows might be nil
			if row == nil {
				continue
			}
			var rowKeys []tgbotapi.KeyboardButton
			for _, cmd := range row {
				rowKeys = append(rowKeys, tgbotapi.NewKeyboardButton(string(cmd)))
			}
			keyboard.Keyboard = append(keyboard.Keyboard, rowKeys)
		}

		msg.ReplyMarkup = keyboard

	} else if len(options.inlineKeyboard) > 0 {

		markup := tgbotapi.NewInlineKeyboardMarkup()
		for _, row := range options.inlineKeyboard {
			keyboardRow := tgbotapi.NewInlineKeyboardRow()
			for _, button := range row {
				keyboardRow = append(keyboardRow, tgbotapi.NewInlineKeyboardButtonData(button.Label, button.Data))
			}
			markup.InlineKeyboard = append(markup.InlineKeyboard, keyboardRow)
		}
		msg.ReplyMarkup = markup
	} else {
		if !options.keepKeyboard {
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		}
	}
	msg.DisableNotification = !options.notification

	sentMsg, err := bc.botApi.Send(msg)
	if err != nil {
		log.Printf("Error sending message %#v: %v", msg, err)
	}
	return sentMsg.MessageID
}

type (
	sendMessageOptions struct {
		keepKeyboard   bool
		inlineKeyboard InlineKeyboard
		notification   bool
	}
	SendMessageOption func(options *sendMessageOptions)
)

func SendMessageKeepKeyboard() SendMessageOption {
	return func(opts *sendMessageOptions) {
		opts.keepKeyboard = true
	}
}

func SendMessageInlineKeyboard(keyboard InlineKeyboard) SendMessageOption {
	return func(opts *sendMessageOptions) {
		opts.inlineKeyboard = keyboard
	}
}

func SendMessageWithNotification() SendMessageOption {
	return func(opts *sendMessageOptions) {
		opts.notification = true
	}
}

func (bc *BotSession) SendMessage(text string, opts ...SendMessageOption) int {
	return bc.SendMessageWithCommands(text, nil, opts...)
}

func (bc *BotSession) UpdateMessageForCallback(queryId string, messageId int, text string, opts ...SendMessageOption) {
	edit := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    bc.chatId,
			MessageID: messageId,
		},
		Text:      text,
		ParseMode: "html",
	}

	options := &sendMessageOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if len(options.inlineKeyboard) > 0 {
		edit.BaseEdit.ReplyMarkup = convertToMarkup(options.inlineKeyboard)
	}

	_, err := bc.botApi.Request(edit)
	if err != nil {
		log.Printf("error updating message: %v", err)
	}
	bc.botApi.Request(tgbotapi.NewCallback(queryId, ""))
}

func (bc *BotSession) SendError(err error) {
	_, sendErr := bc.botApi.Send(tgbotapi.NewMessage(bc.ChatId(), fmt.Sprintf("error: %v", err)))
	if sendErr != nil {
		log.Printf("Error sending error: %v", sendErr)
	}
}

func (bc *BotSession) Fail(message string, formatErrorMsg string, args ...interface{}) {
	log.Printf(formatErrorMsg, args...)
	bc.SendMessage(message)
	bc.PopState()
}

func (bc *BotSession) AcceptUsers() {
	bc.bot.AcceptUsers(bc.botCtx)
}

func (bc *BotSession) BotName() (string, error) {
	me, err := bc.botApi.GetMe()
	if err != nil {
		return "", fmt.Errorf("error getting bot identity: %v", err)
	}
	return me.UserName, nil
}

func (bc *BotSession) Shutdown() {
	for i := len(bc.states) - 1; i >= 0; i-- {
		bc.states[i].BeforeLeave(bc)
	}
}

func convertToMarkup(keyboard InlineKeyboard) *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for _, row := range keyboard {
		keyboardRow := tgbotapi.NewInlineKeyboardRow()
		for _, button := range row {
			keyboardRow = append(keyboardRow, tgbotapi.NewInlineKeyboardButtonData(button.Label, button.Data))
		}
		markup.InlineKeyboard = append(markup.InlineKeyboard, keyboardRow)
	}
	return &markup
}
