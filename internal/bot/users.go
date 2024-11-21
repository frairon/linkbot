package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/frairon/botty"
	"github.com/frairon/linkbot/internal/storage/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UsersList() State {
	var (
		Add    Button = "➕ Add"
		Back   Button = "↩ Back"
		Delete Button = "❌ Delete"
	)

	var users []*models.User

	return &functionState{
		activate: func(bs *BotSession) {
			var err error
			users, err = bs.Storage().ListUsers()
			if err != nil {
				bs.Fail("Cannot list users", "error reading users: %v", err)
				return
			}

			content, err := botty.RunTemplate(`All Users
{{divider}}
{{- if .users -}}
{{- range $idx, $user:= .users }}
[{{$idx}}] {{$user.Name.String}} ({{$user.ID}})
{{- end -}}
{{- else }}
- no users registered -
{{- end -}}`, botty.KV("users", users))

			if err != nil {
				bs.Fail("Cannot list users", "error reading users: %v", err)
				return
			} else {
				bs.SendMessageWithCommands(content, NewButtonKeyboard(newRow(Back),
					newRow(Add, Delete)))
			}
		},
		handleMessage: func(bs *BotSession, message *tgbotapi.Message) {
			botName, err := bs.BotName()
			if err != nil {
				bs.Fail("Cannot find bot identity", "error getting bot name: %v", err)
				return
			}

			switch Button(message.Text) {
			case Back:
				bs.PopState()
			case Add:
				text, err := botty.RunTemplate(`The bot is now set to ACCEPT-mode, allowing new users to join.
This will be disabled automatically after 10 minutes.
Tell you friend to contact bot @{{.botName}} now.`, botty.KV("botName", botName))
				if err != nil {
					bs.Fail("error rendering template", "error rendering template: %v", err)
					return
				}
				bs.SendMessageWithCommands(text, nil)
				bs.AcceptUsers()
			case Delete:
				bs.PushState(SelectToDeleteUser(users))
			}
		},
	}
}

func SelectToDeleteUser(users []*models.User) State {
	var Back Button = "Back"
	return &functionState{
		activate: func(bs *BotSession) {
			bs.SendMessageWithCommands("Select user to delete", NewButtonKeyboard(newRow(Back)))
		},
		handleMessage: func(bs *BotSession, msg *tgbotapi.Message) {
			selector := strings.TrimSpace(msg.Text)

			idx, err := strconv.ParseInt(selector, 10, 32)
			if err != nil || idx < 0 || int(idx) >= len(users) {
				bs.SendMessage(fmt.Sprintf("Cannot find user by '%s'. Enter valid index.", selector))
				return
			}

			user := users[idx]

			bs.ReplaceState(PromptState(func() {
				err := bs.Storage().DeleteUser(user)
				if err != nil {
					log.Printf("error deleting item %#v: %v", user, err)
					bs.SendMessage("error deleting user")
				}
			}))
		},
	}
}
