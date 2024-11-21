package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Button         string
	ButtonRow      []Button
	ButtonKeyboard []ButtonRow
)

func (c Button) Is(val string) bool {
	return string(c) == val
}

func (c Button) S() string {
	return string(c)
}

type State interface {
	Activate(bs *BotSession)
	Return(bs *BotSession)
	HandleMessage(bs *BotSession, message *tgbotapi.Message) bool
	HandleCommand(bs *BotSession, command string, args ...string) bool
	HandleCallbackQuery(bs *BotSession, query *tgbotapi.CallbackQuery) bool

	// called before leaving the state (either by pushing another state on top of it or popping it)
	BeforeLeave(bs *BotSession)
}

type StateFactory func() State

func NewButtonKeyboard(rows ...ButtonRow) ButtonKeyboard {
	return ButtonKeyboard(rows)
}

func newConditionalRow(condition func() bool, row ButtonRow) ButtonRow {
	if condition() {
		return row
	}
	return nil
}

func newRow(commands ...Button) ButtonRow {
	return ButtonRow(commands)
}

func conditionalButton(condition func() bool, trueButton, falseButton Button) Button {
	if condition() {
		return trueButton
	}
	return falseButton
}

type (
	InlineButton struct {
		Label string
		Data  string
	}
	InlineRow      []InlineButton
	InlineKeyboard []InlineRow
)

func NewInlineKeyboard(rows ...InlineRow) InlineKeyboard {
	return InlineKeyboard(rows)
}

func NewInlineRow(buttons ...InlineButton) InlineRow {
	return InlineRow(buttons)
}

func NewInlineButton(label, data string) InlineButton {
	return InlineButton{
		Label: label,
		Data:  data,
	}
}

type functionState struct {
	activate             func(bs *BotSession)
	returner             func(bs *BotSession)
	handleMessage        func(bs *BotSession, message *tgbotapi.Message)
	commandHandler       func(bs *BotSession, command string, args ...string) bool
	callbackQueryHandler func(bs *BotSession, query *tgbotapi.CallbackQuery) bool
	beforeLeaveHandler   func(bs *BotSession)
}

type StateBuilder struct {
	fs *functionState
}

func NewStateBuilder() *StateBuilder {
	return &StateBuilder{
		fs: &functionState{
			activate: func(bs *BotSession) {
				bs.SendMessage("I am a state")
			},
		},
	}
}

func (sb *StateBuilder) OnActivate(activator func(bs *BotSession)) *StateBuilder {
	sb.fs.activate = activator
	return sb
}
func (sb *StateBuilder) Build() State {
	return sb.fs
}

func (fs *functionState) Activate(bc *BotSession) {
	fs.activate(bc)
}

func (fs *functionState) Return(bs *BotSession) {
	if fs.returner != nil {
		fs.returner(bs)
	} else {
		fs.activate(bs)
	}
}

func (fs *functionState) HandleMessage(bs *BotSession, message *tgbotapi.Message) bool {
	if fs.handleMessage == nil {
		return false
	}
	fs.handleMessage(bs, message)
	return true
}

func (fs *functionState) HandleCommand(bs *BotSession, command string, args ...string) bool {
	if fs.commandHandler != nil {
		return fs.commandHandler(bs, command, args...)
	}
	return false
}

func (fs *functionState) HandleCallbackQuery(bs *BotSession, query *tgbotapi.CallbackQuery) bool {
	if fs.callbackQueryHandler != nil {
		return fs.callbackQueryHandler(bs, query)
	}
	return false
}

func (fs *functionState) BeforeLeave(bs *BotSession) {
	if fs.beforeLeaveHandler != nil {
		fs.beforeLeaveHandler(bs)
	}
}

const (
	HomeButton Button = "🏠"
)

func GlobalHomeHandler(bs *BotSession, message *tgbotapi.Message) bool {
	if message.Text == HomeButton.S() {
		bs.ResetToState(bs.bot.RootState())
		return true
	}
	return false
}
