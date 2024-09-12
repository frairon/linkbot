package link

import "github.com/frairon/linkbot/internal/bot"

func Home() bot.State {
	return bot.NewStateBuilder().Done()
}
