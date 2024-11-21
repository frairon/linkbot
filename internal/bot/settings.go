package bot

import (
	"log"

	"github.com/frairon/botty"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ConfigureSessionSettings() State {
	var (
		BellOn          Button = "Bell 🔔"
		BellOff         Button = "Bell 🔕"
		ButtonActionOn  Button = "Button Action 🔔"
		ButtonActionOff Button = "Button Action 🔕"
		AutomationOn    Button = "Automation 🔔"
		AutomationOff   Button = "Automation 🔕"
		MovementOn      Button = "Movement 🔔"
		MovementOff     Button = "Movement 🔕"
		EnterLeaveOn    Button = "Enter/Leave 🔔"
		EnterLeaveOff   Button = "Enter/Leave 🔕"
		ExternalOn      Button = "External 🔔"
		ExternalOff     Button = "External 🔕"
		MuteAll         Button = "Mute All"
		Unmute          Button = "Unmute"
	)
	return &functionState{
		activate: func(bs *BotSession) {
			content, err := botty.RunTemplate(`<b>User Settings</b>
========
{{if .settings.PauseAllNotifications -}}
All Notifications Paused
{{- else}}
Notify on:
  • Bell {{if .settings.NotifyOnBell}}🔔{{else}}🔕{{end}}
  • Button Action {{if .settings.NotifyOnButtonTempChange}}🔔{{else}}🔕{{end}}
  • Automation Action  {{if .settings.NotifyOnAutoTempChange}}🔔{{else}}🔕{{end}}
  • Movement {{if .settings.NotifyOnMovement}}🔔{{else}}🔕{{end}}
  • External Action  {{if .settings.NotifyOnTempChange}}🔔{{else}}🔕{{end}}
  • Enter/Leave of device  {{if .settings.NotifyOnEnterLeave}}🔔{{else}}🔕{{end}}
{{end}}`,
				botty.KV("settings", bs.settings))
			if err != nil {
				bs.SendError(err)
			}
			bs.SendMessageWithCommands(content,
				NewButtonKeyboard(
					newRow(Back, conditionalButton(func() bool { return bs.settings.PauseAllNotifications }, Unmute, MuteAll)),
					newConditionalRow(func() bool { return !bs.settings.PauseAllNotifications },
						newRow(
							conditionalButton(func() bool { return bs.settings.NotifyOnBell }, BellOff, BellOn),
							conditionalButton(func() bool { return bs.settings.NotifyOnButtonTempChange }, ButtonActionOff, ButtonActionOn),
						),
					),
					newConditionalRow(func() bool { return !bs.settings.PauseAllNotifications },
						newRow(
							conditionalButton(func() bool { return bs.settings.NotifyOnAutoTempChange }, AutomationOff, AutomationOn),
							conditionalButton(func() bool { return bs.settings.NotifyOnMovement }, MovementOff, MovementOn),
							conditionalButton(func() bool { return bs.settings.NotifyOnTempChange }, ExternalOff, ExternalOn),
							conditionalButton(func() bool { return bs.settings.NotifyOnEnterLeave }, EnterLeaveOff, EnterLeaveOn),
						),
					),
				))
		},
		handleMessage: func(bs *BotSession, message *tgbotapi.Message) {
			switch Button(message.Text) {
			case Back:
				bs.PopState()
				return
			case BellOn:
				bs.settings.NotifyOnBell = true
			case BellOff:
				bs.settings.NotifyOnBell = false
			case ButtonActionOn:
				bs.settings.NotifyOnButtonTempChange = true
			case ButtonActionOff:
				bs.settings.NotifyOnButtonTempChange = false
			case AutomationOn:
				bs.settings.NotifyOnAutoTempChange = true
			case AutomationOff:
				bs.settings.NotifyOnAutoTempChange = false
			case MovementOn:
				bs.settings.NotifyOnMovement = true
			case MovementOff:
				bs.settings.NotifyOnMovement = false
			case EnterLeaveOn:
				bs.settings.NotifyOnEnterLeave = true
			case EnterLeaveOff:
				bs.settings.NotifyOnEnterLeave = false
			case ExternalOn:
				bs.settings.NotifyOnTempChange = true
			case ExternalOff:
				bs.settings.NotifyOnTempChange = false
			case MuteAll:
				bs.settings.PauseAllNotifications = true
			case Unmute:
				bs.settings.PauseAllNotifications = false
			}
			bs.CurrentState().Activate(bs)
		},
	}
}

func ConfigureMultiSettings() State {
	renderSettings := func(name string, value bool) (string, error) {
		return botty.RunTemplate(`{{.name}} {{if .value}}🗹{{else}}☐{{end}}`,
			botty.KV("name", name), botty.KV("value", value))
	}

	var (
		Disable = NewInlineButton("Disable", "disable")
		Enable  = NewInlineButton("Enable", "enable")
	)

	makeKeyboard := func(value bool) InlineKeyboard {
		return NewInlineKeyboard(NewInlineRow(TernaryButton(value, Disable, Enable)))
	}

	return NewMultiMessageHandler(
		func(bs *BotSession, query string) (string, InlineKeyboard, error) {
			switch query {
			case Disable.Data:
				bs.settings.PauseAllNotifications = false
			case Enable.Data:
				bs.settings.PauseAllNotifications = true
			default:
				log.Printf("unhandled query: %s", query)
			}

			content, err := renderSettings("Pause all Notifications", bs.settings.PauseAllNotifications)
			return content, makeKeyboard(bs.settings.PauseAllNotifications), err
		},
		func(bs *BotSession, query string) (string, InlineKeyboard, error) {
			switch query {
			case Disable.Data:
				bs.settings.NotifyOnAutoTempChange = false
			case Enable.Data:
				bs.settings.NotifyOnAutoTempChange = true
			default:
				log.Printf("unhandled query: %s", query)
			}

			content, err := renderSettings("Notify on Temp Change", bs.settings.NotifyOnAutoTempChange)
			return content, makeKeyboard(bs.settings.NotifyOnAutoTempChange), err
		},
		func(bs *BotSession, query string) (string, InlineKeyboard, error) {
			switch query {
			case Disable.Data:
				bs.settings.NotifyOnBell = false
			case Enable.Data:
				bs.settings.NotifyOnBell = true
			default:
				log.Printf("unhandled query: %s", query)
			}

			content, err := renderSettings("Notify on Bell", bs.settings.NotifyOnBell)
			return content, makeKeyboard(bs.settings.NotifyOnBell), err
		},
	)
}

func TernaryButton(cond bool, trueButton, falseButton InlineButton) InlineButton {
	if cond {
		return trueButton
	}
	return falseButton
}

func Settings() State {
	renderSettings := func(settings *SessionSettings) (string, error) {
		return botty.RunTemplate(`<b>User Settings</b>
		========
		PauseAllNotifications  {{if .settings.PauseAllNotifications}}🗹{{else}}☐{{end}}
		`,
			botty.KV("settings", settings))
	}

	var (
		PauseAllNotificationsOn  = NewInlineButton("🔇 All Notifications", "all_notifications_off")
		PauseAllNotificationsOff = NewInlineButton("🔈 All Notifications", "all_notifications_on")
	)

	makeKeyboard := func(settings *SessionSettings) InlineKeyboard {
		exitRow := NewInlineRow(NewInlineButton(Back.S(), "back"))

		return NewInlineKeyboard(
			NewInlineRow(TernaryButton(settings.PauseAllNotifications, PauseAllNotificationsOff, PauseAllNotificationsOn)),
			exitRow,
		)
	}

	return NewMessageHandler(func(bs *BotSession, query string) (string, InlineKeyboard, error) {
		switch query {
		case PauseAllNotificationsOff.Data:
			log.Printf("pausing notifications")
			bs.settings.PauseAllNotifications = false
		case PauseAllNotificationsOn.Data:
			bs.settings.PauseAllNotifications = true
		case "back":
			bs.PopState()
			return "", nil, nil
		default:
			log.Printf("unhandled query: %s", query)
		}

		content, err := renderSettings(bs.settings)
		return content, makeKeyboard(bs.settings), err
	})
}

type InlineMessageHandler func(bs *BotSession, query string) (string, InlineKeyboard, error)

func NewMessageHandler(handleQuery InlineMessageHandler) State {
	var lastMessageId int

	fs := &functionState{
		activate: func(bs *BotSession) {
			msg, keyboard, err := handleQuery(bs, "")
			if err != nil {
				bs.SendError(err)
				return
			}
			lastMessageId = bs.SendMessage(msg, SendMessageInlineKeyboard(keyboard))
		},
		callbackQueryHandler: func(bs *BotSession, query *tgbotapi.CallbackQuery) bool {
			log.Printf("callback: %#v", query)
			content, keyboard, err := handleQuery(bs, query.Data)
			if err != nil {
				bs.SendError(err)
				return true
			}
			if content != "" && keyboard != nil {
				bs.UpdateMessageForCallback(query.ID,
					query.Message.MessageID,
					content,
					SendMessageInlineKeyboard(keyboard),
				)
			}
			return true
		},
		beforeLeaveHandler: func(bs *BotSession) {
			if lastMessageId != 0 {
				bs.RemoveKeyboardForMessage(lastMessageId)
			}
		},
	}
	return fs
}

func NewMultiMessageHandler(handlers ...InlineMessageHandler) State {
	handlersByMsg := map[int]InlineMessageHandler{}

	fs := &functionState{
		activate: func(bs *BotSession) {
			for _, handler := range handlers {
				msg, keyboard, err := handler(bs, "")
				if err != nil {
					bs.SendError(err)
					return
				}
				msgId := bs.SendMessage(msg, SendMessageInlineKeyboard(keyboard))
				handlersByMsg[msgId] = handler
			}
		},
		callbackQueryHandler: func(bs *BotSession, query *tgbotapi.CallbackQuery) bool {
			handler := handlersByMsg[query.Message.MessageID]

			if handler == nil {
				log.Printf("did not find handler for message")
				return false
			}
			content, keyboard, err := handler(bs, query.Data)
			if err != nil {
				bs.SendError(err)
				return true
			}
			if content != "" && keyboard != nil {
				bs.UpdateMessageForCallback(query.ID,
					query.Message.MessageID,
					content,
					SendMessageInlineKeyboard(keyboard),
				)
			}
			return true
		},
		beforeLeaveHandler: func(bs *BotSession) {
			for msgId := range handlersByMsg {
				bs.RemoveKeyboardForMessage(msgId)
			}
		},
	}
	return fs
}
