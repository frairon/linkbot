package link

import (
	"errors"
	"fmt"

	"github.com/frairon/botty"
)

func Home() LinkState {

	return botty.NewStateBuilder[*State]().OnActivate(func(bs LinkSession) {

		type categoryView = struct {
			Title string
			Count int
		}

		catsModel, err := bs.State().st.ListCategories(int64(bs.UserId()))

		var cats []categoryView
		for _, cat := range catsModel {
			cats = append(cats, categoryView{
				Title: cat.Category,
				Count: cat.Count,
			})
		}

		bs.SendError(err)

		bs.SendTemplateMessage(`
Categories:
{{- range $cat := .cats }}
- {{$cat.Title}} ({{$cat.Count}})
{{end }}
		`, botty.TplValues(botty.KV("cats", cats)),
			botty.SendMessageWithKeyboard(botty.NewButtonKeyboard(
				botty.ButtonRow{
					"back", "add",
				},
			),
			))
	}).OnMessage(func(bs botty.Session[*State], message botty.ChatMessage) {
		switch message.Text() {
		case "back":
			bs.PopState()
		case "add":
			bs.PushState(AddLink())
		}
	}).
		Build()
}

func AddLink() LinkState {

	render := func(bs LinkSession) {
		bs.SendMessage("Paste the link and send it to me please")
	}

	var li *Link

	return botty.NewStateBuilder[*State]().OnActivate(func(bs LinkSession) {
		render(bs)
	}).
		OnMessage(func(bs botty.Session[*State], message botty.ChatMessage) {
			if li != nil {
				switch message.Text() {
				case "Discard":
					bs.SendMessage("Discarding.")
					bs.PopState()
					return
				case "Add":
					bs.State().st.AddLink(int64(bs.UserId()), "somecategory", li.Url, li.Headline)
					bs.SendMessage("Added link.")
					bs.PopState()
					return
				default:
					bs.SendError(fmt.Errorf("unhandled message. Canceling operation"))
					bs.PopState()
					return
				}
			}
			var err error
			li, err = parseLink(message.Text())
			if err != nil {
				bs.SendMessage(fmt.Sprintf("Error parsing the link: %s", errors.Unwrap(err).Error()))
				bs.SendMessage("You want to try again?")
				return
			}
			bs.SendTemplateMessage(`
Received Link {{.li.Url}}.
Extracted headline: {{.li.Headline}}
`, botty.TplValues(botty.KV("li", li)),
				botty.SendMessageWithKeyboard(botty.NewButtonKeyboard(botty.ButtonRow{
					"Add", "Discard",
				})),
			)
		}).
		Build()
}
