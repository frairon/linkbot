package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type promptOptions struct {
	dropStates int
	message    string
}

type PromptOption func(opts *promptOptions)

func PromptDropStates(states int) PromptOption {
	return func(opts *promptOptions) {
		opts.dropStates = states
	}
}

func PromptMessagef(format string, args ...interface{}) PromptOption {
	return func(opts *promptOptions) {
		opts.message = fmt.Sprintf(format, args...)
	}
}

func PromptState(yesHandler func(), options ...PromptOption) State {
	const (
		Yes    Button = "⚠ Yes"
		Cancel Button = "Cancel"
	)

	opts := &promptOptions{
		dropStates: 1,
		message:    "Are you sure?",
	}

	for _, option := range options {
		option(opts)
	}

	return &functionState{
		activate: func(bs *botSession) {
			bs.SendMessageWithCommands(opts.message, NewButtonKeyboard(newRow(Yes, Cancel)))
		},

		handleMessage: func(bs *botSession, message *tgbotapi.Message) {
			switch Button(message.Text) {
			case Cancel:
				bs.SendMessage("Aborted.")
				bs.DropStates(opts.dropStates)
			case Yes:
				yesHandler()
				bs.DropStates(opts.dropStates)
			}
		},
	}
}

func SelectState[T any](text string, items []T, accept func(bs *botSession, item T)) State {
	return &functionState{
		activate: func(bs *botSession) {
			bs.SendMessage(text)
			bs.SendMessage(fmt.Sprintf("Please enter index (0-%d)", len(items)-1))
		},
		handleMessage: func(bs *botSession, msg *tgbotapi.Message) {
			selector := strings.TrimSpace(msg.Text)

			idx, err := strconv.ParseInt(selector, 10, 32)
			if err != nil || idx < 0 || int(idx) >= len(items) {
				bs.SendMessage(fmt.Sprintf("Cannot find Item by '%s'. Enter valid item.", selector))
				return
			}

			accept(bs, items[idx])
			bs.PopState()
		},
	}
}

const (
	NewList  Button = "✳ New List"
	Add      Button = "➕"
	AddMulti Button = "➕➕"
	Delete   Button = "❌ Delete"

	ListEditMode Button = "🗂"
	Select       Button = "Select"

	Back Button = "🔙"
	More Button = "..."

	Rename  Button = "✏ Rename"
	Sharing Button = "👥 Sharing ..."

	Share   Button = "➕👥 Share"
	Unshare Button = "❌👥 Unshare"
)
