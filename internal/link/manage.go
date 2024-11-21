package link

import (
	"github.com/frairon/botty"
)

func Home() LinkState {
	return botty.NewStateBuilder[*State]().OnActivate(func(bs LinkSession) {
		bs.SendTemplateMessage("asd", botty.TplValues(botty.KV("key", "value")), botty.SendMessageWithKeyboard(botty.NewButtonKeyboard(botty.ButtonRow{
			"back",
		})))
	}).Build()
}
