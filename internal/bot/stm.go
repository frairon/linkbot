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
	Activate(bs *botSession)
	Return(bs *botSession)
	HandleMessage(bs *botSession, message *tgbotapi.Message) bool
	HandleCommand(bs *botSession, command string, args ...string) bool
	HandleCallbackQuery(bs *botSession, query *tgbotapi.CallbackQuery) bool

	// called before leaving the state (either by pushing another state on top of it or popping it)
	BeforeLeave(bs *botSession)
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
	activate             func(bs *botSession)
	returner             func(bs *botSession)
	handleMessage        func(bs *botSession, message *tgbotapi.Message)
	commandHandler       func(bs *botSession, command string, args ...string) bool
	callbackQueryHandler func(bs *botSession, query *tgbotapi.CallbackQuery) bool
	beforeLeaveHandler   func(bs *botSession)
}

type StateBuilder struct {
	fs *functionState
}

func NewStateBuilder() *StateBuilder {
	return &StateBuilder{
		fs: &functionState{
			activate: func(bs *botSession) {
				bs.SendMessage("I am a state")
			},
		},
	}
}

func (sb *StateBuilder) Done() State {
	return sb.fs
}

func (fs *functionState) Activate(bc *botSession) {
	fs.activate(bc)
}

func (fs *functionState) Return(bs *botSession) {
	if fs.returner != nil {
		fs.returner(bs)
	} else {
		fs.activate(bs)
	}
}

func (fs *functionState) HandleMessage(bs *botSession, message *tgbotapi.Message) bool {
	if fs.handleMessage == nil {
		return false
	}
	fs.handleMessage(bs, message)
	return true
}

func (fs *functionState) HandleCommand(bs *botSession, command string, args ...string) bool {
	if fs.commandHandler != nil {
		return fs.commandHandler(bs, command, args...)
	}
	return false
}

func (fs *functionState) HandleCallbackQuery(bs *botSession, query *tgbotapi.CallbackQuery) bool {
	if fs.callbackQueryHandler != nil {
		return fs.callbackQueryHandler(bs, query)
	}
	return false
}

func (fs *functionState) BeforeLeave(bs *botSession) {
	if fs.beforeLeaveHandler != nil {
		fs.beforeLeaveHandler(bs)
	}
}

const (
	HomeButton Button = "🏠"
)

func GlobalHomeHandler(bs *botSession, message *tgbotapi.Message) bool {
	if message.Text == HomeButton.S() {
		bs.ResetToState(bs.bot.RootState())
		return true
	}
	return false
}
